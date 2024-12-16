package types

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type PKSigner interface {
	GetPublicKey(tx *Transaction) ([]byte, error)
}

func NewPKSigner(chainId *big.Int) PKSigner {
	return cancunSigner{londonSigner{eip2930Signer{NewEIP155Signer(chainId)}}}
}

func (s cancunSigner) GetPublicKey(tx *Transaction) ([]byte, error) {
	if tx.Type() != BlobTxType {
		return s.londonSigner.GetPublicKey(tx)
	}
	V, R, S := tx.RawSignatureValues()
	// Blob txs are defined to use 0 and 1 as their recovery
	// id, add 27 to become equivalent to unprotected Homestead signatures.
	V = new(big.Int).Add(V, big.NewInt(27))
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, fmt.Errorf("%w: have %d want %d", ErrInvalidChainId, tx.ChainId(), s.chainId)
	}
	return recoverPlainForPubK(s.Hash(tx), R, S, V, true)
}

func (s londonSigner) GetPublicKey(tx *Transaction) ([]byte, error) {
	if tx.Type() != DynamicFeeTxType {
		return s.eip2930Signer.GetPublicKey(tx)
	}
	V, R, S := tx.RawSignatureValues()
	// DynamicFee txs are defined to use 0 and 1 as their recovery
	// id, add 27 to become equivalent to unprotected Homestead signatures.
	V = new(big.Int).Add(V, big.NewInt(27))
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, fmt.Errorf("%w: have %d want %d", ErrInvalidChainId, tx.ChainId(), s.chainId)
	}
	return recoverPlainForPubK(s.Hash(tx), R, S, V, true)
}

func (s eip2930Signer) GetPublicKey(tx *Transaction) ([]byte, error) {
	V, R, S := tx.RawSignatureValues()
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.GetPublicKey(tx)
	case AccessListTxType:
		// AL txs are defined to use 0 and 1 as their recovery
		// id, add 27 to become equivalent to unprotected Homestead signatures.
		V = new(big.Int).Add(V, big.NewInt(27))
	default:
		return nil, ErrTxTypeNotSupported
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, fmt.Errorf("%w: have %d want %d", ErrInvalidChainId, tx.ChainId(), s.chainId)
	}
	return recoverPlainForPubK(s.Hash(tx), R, S, V, true)
}

func (s EIP155Signer) GetPublicKey(tx *Transaction) ([]byte, error) {
	if tx.Type() != LegacyTxType {
		return nil, ErrTxTypeNotSupported
	}
	if !tx.Protected() {
		return HomesteadSigner{}.GetPublicKey(tx)
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, fmt.Errorf("%w: have %d want %d", ErrInvalidChainId, tx.ChainId(), s.chainId)
	}
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlainForPubK(s.Hash(tx), R, S, V, true)
}

func (hs HomesteadSigner) GetPublicKey(tx *Transaction) ([]byte, error) {
	if tx.Type() != LegacyTxType {
		return nil, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlainForPubK(hs.Hash(tx), r, s, v, true)
}

func recoverPlainForPubK(sighash common.Hash, R, S, Vb *big.Int, homestead bool) ([]byte, error) {
	if Vb.BitLen() > 8 {
		return nil, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return nil, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return nil, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return nil, errors.New("invalid public key")
	}
	return pub, nil
}
