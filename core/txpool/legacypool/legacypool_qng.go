package legacypool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (pool *LegacyPool) RemoveTx(hash common.Hash, outofbound bool) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.removeTx(hash, outofbound, true)
}

func (pool *LegacyPool) All() *lookup {
	return pool.all
}

func (pool *LegacyPool) AddRemotesSync(txs []*types.Transaction) []error {
	return pool.addRemotesSync(txs)
}
