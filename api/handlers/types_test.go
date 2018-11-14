/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"testing"

	"github.com/seeleteam/scan-api/database"
	"github.com/stretchr/testify/assert"
)

func newTestDBBlock(t *testing.T) *database.DBBlock {
	return &database.DBBlock{
		Reward:  99,
		Height:  133853,
		UsedGas: 900,
		Txs:     []database.DBSimpleTxInBlock{{Fee: 99, GasPrice: 88}},
	}
}

func Test_CreateRetSimpleBlockInfo(t *testing.T) {
	var TxFee, TxPrice int64
	header := newTestDBBlock(t)
	got := createRetSimpleBlockInfo(header)

	for i := 0; i < len(header.Txs); i++ {
		TxFee += header.Txs[i].Fee
		TxPrice += header.Txs[i].GasPrice
	}
	assert.Equal(t, got.Fee, TxFee)
	assert.Equal(t, got.Gasprice, TxPrice)
	assert.Equal(t, got.UsedGas, header.UsedGas)
	assert.Equal(t, int64(got.Height), header.Height)
	assert.Equal(t, got.Reward, header.Reward)
}

func Test_CreateRetLastblockInfo(t *testing.T) {
	var lastblockHeight, lastblockTime int64
	lastblockHeight = 10399
	lastblockTime = 12
	got := createRetLastblockInfo(lastblockHeight, lastblockTime)

	assert.Equal(t, got.LastblockHeight, lastblockHeight)
	assert.Equal(t, got.LastblockTime, lastblockTime)
}

func Benchmark_isOneBitCharacter(b *testing.B) {
	tests := struct {
		lastblockHeight int64
		lastblockTime   int64
	}{10399, 12}

	for idx := 0; idx < b.N; idx++ {
		if got := createRetLastblockInfo(tests.lastblockHeight, tests.lastblockTime); got.LastblockHeight != tests.lastblockHeight {
			b.Errorf("createRetLastblockInfo() = %v, lastblockHeight %v", got, tests.lastblockHeight)
		}
	}

}
