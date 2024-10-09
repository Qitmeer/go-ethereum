package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGetPubkeyFromTx(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := NewEIP155Signer(big.NewInt(18))
	tx, err := SignTx(NewTransaction(0, addr, new(big.Int), 0, new(big.Int), nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if from != addr {
		t.Errorf("expected from and address to be equal. Got %x want %x", from, addr)
	}

	var gs PKSigner
	gs = &signer
	pkb, err := gs.GetPublicKey(tx)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := crypto.UnmarshalPubkey(pkb)
	if err != nil {
		t.Fatal(err)
	}
	pka := crypto.PubkeyToAddress(*pk)
	if pka != addr {
		t.Errorf("expected pubkey-address and address to be equal. Got %x want %x", pka, addr)
	}
}
