package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"math/big"
)

type IMiner interface {
	SetExtra(extra []byte) error
	Pending() (*types.Block, types.Receipts, *state.StateDB)
	SetGasCeil(ceil uint64)
	SetGasTip(tip *big.Int) error
	BuildPayload(args *BuildPayloadArgs) (*Payload, error)
}

type TransactionsByPriceAndNonce interface {
	Peek() (*txpool.LazyTransaction, *uint256.Int)
	Shift()
	Pop()
}

func NewTransactionsByPriceAndNonce(signer types.Signer, txs map[common.Address][]*txpool.LazyTransaction, baseFee *big.Int) TransactionsByPriceAndNonce {
	return newTransactionsByPriceAndNonce(signer, txs, baseFee)
}
