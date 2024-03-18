package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/holiman/uint256"
	"math/big"
	"time"
)

type IMiner interface {
	Start()
	Stop()
	Close()
	Mining() bool
	Hashrate() uint64
	SetExtra(extra []byte) error
	SetRecommitInterval(interval time.Duration)
	Pending() (*types.Block, *state.StateDB)
	PendingBlock() *types.Block
	PendingBlockAndReceipts() (*types.Block, types.Receipts)
	SetEtherbase(addr common.Address)
	SetGasCeil(ceil uint64)
	SetGasTip(tip *big.Int) error
	SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription
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
