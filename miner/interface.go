package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"time"
)

type IMiner interface {
	Start(coinbase common.Address)
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
	EnablePreseal()
	DisablePreseal()
	SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription
	GetSealingBlockAsync(parent common.Hash, timestamp uint64, coinbase common.Address, random common.Hash, noTxs bool) (chan *types.Block, error)
	GetSealingBlockSync(parent common.Hash, timestamp uint64, coinbase common.Address, random common.Hash, noTxs bool) (*types.Block, error)
}
