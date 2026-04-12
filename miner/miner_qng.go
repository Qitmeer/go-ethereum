package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

func (miner *Miner) ForcePending() (*types.Block, types.Receipts, *state.StateDB) {
	miner.pending.update(common.Hash{}, nil)
	return miner.Pending()
}
