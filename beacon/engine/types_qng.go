package engine

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"math/big"
)

func ExecutableDataToBlockQng(params ExecutableData, versionedHashes []common.Hash, beaconRoot *common.Hash) (*types.Block, error) {
	txs, err := decodeTransactions(params.Transactions)
	if err != nil {
		return nil, err
	}
	if len(params.LogsBloom) != 256 {
		return nil, fmt.Errorf("invalid logsBloom length: %v", len(params.LogsBloom))
	}
	// Check that baseFeePerGas is not negative or too big
	if params.BaseFeePerGas != nil && (params.BaseFeePerGas.Sign() == -1 || params.BaseFeePerGas.BitLen() > 256) {
		return nil, fmt.Errorf("invalid baseFeePerGas: %v", params.BaseFeePerGas)
	}
	var blobHashes []common.Hash
	for _, tx := range txs {
		blobHashes = append(blobHashes, tx.BlobHashes()...)
	}
	if len(blobHashes) != len(versionedHashes) {
		return nil, fmt.Errorf("invalid number of versionedHashes: %v blobHashes: %v", versionedHashes, blobHashes)
	}
	for i := 0; i < len(blobHashes); i++ {
		if blobHashes[i] != versionedHashes[i] {
			return nil, fmt.Errorf("invalid versionedHash at %v: %v blobHashes: %v", i, versionedHashes, blobHashes)
		}
	}
	// Only set withdrawalsRoot if it is non-nil. This allows CLs to use
	// ExecutableData before withdrawals are enabled by marshaling
	// Withdrawals as the json null value.
	var withdrawalsRoot *common.Hash
	if params.Withdrawals != nil {
		h := types.DeriveSha(types.Withdrawals(params.Withdrawals), trie.NewStackTrie(nil))
		withdrawalsRoot = &h
	}
	header := &types.Header{
		ParentHash:       params.ParentHash,
		UncleHash:        types.EmptyUncleHash,
		Coinbase:         params.FeeRecipient,
		Root:             params.StateRoot,
		TxHash:           types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil)),
		ReceiptHash:      params.ReceiptsRoot,
		Bloom:            types.BytesToBloom(params.LogsBloom),
		Difficulty:       params.Difficulty,
		Number:           new(big.Int).SetUint64(params.Number),
		GasLimit:         params.GasLimit,
		GasUsed:          params.GasUsed,
		Time:             params.Timestamp,
		BaseFee:          params.BaseFeePerGas,
		Extra:            params.ExtraData,
		MixDigest:        params.Random,
		WithdrawalsHash:  withdrawalsRoot,
		ExcessBlobGas:    params.ExcessBlobGas,
		BlobGasUsed:      params.BlobGasUsed,
		ParentBeaconRoot: beaconRoot,
	}
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: txs, Uncles: nil, Withdrawals: params.Withdrawals})
	if block.Hash() != params.BlockHash {
		return nil, fmt.Errorf("blockhash mismatch, want %x, got %x", params.BlockHash, block.Hash())
	}
	return block, nil
}
