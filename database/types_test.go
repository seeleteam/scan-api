/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"math/big"
	"testing"

	"github.com/seeleteam/scan-api/rpc"
	"github.com/stretchr/testify/assert"
)

func newTestDBBlock(t *testing.T) *rpc.BlockInfo {
	return &rpc.BlockInfo{
		Hash:      "0x0000",
		Height:    133853,
		Timestamp: new(big.Int),
		Txs:       []rpc.Transaction{{Fee: 99, Amount: new(big.Int)}},
	}
}
func TestCreateDbBlock(t *testing.T) {
	var TxFee int64
	BlockInfo := newTestDBBlock(t)
	got := CreateDbBlock(BlockInfo)

	for i := 0; i < len(BlockInfo.Txs); i++ {
		TxFee += BlockInfo.Txs[i].Fee
	}
	assert.Equal(t, got.Txs[0].Fee, TxFee)
	assert.Equal(t, uint64(got.Height), BlockInfo.Height)
	assert.Equal(t, got.HeadHash, BlockInfo.Hash)
}
