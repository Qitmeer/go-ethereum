package miner

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (payload *Payload) ResolveFullBlock() *types.Block {
	payload.lock.Lock()
	defer payload.lock.Unlock()

	if payload.full == nil {
		select {
		case <-payload.stop:
			return nil
		default:
		}
		// Wait the full payload construction. Note it might block
		// forever if Resolve is called in the meantime which
		// terminates the background construction process.
		payload.cond.Wait()
	}
	// Terminate the background payload construction
	select {
	case <-payload.stop:
	default:
		close(payload.stop)
	}
	return payload.full
}
