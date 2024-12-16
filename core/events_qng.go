package core

import "github.com/ethereum/go-ethereum/core/types"

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *types.Block }
