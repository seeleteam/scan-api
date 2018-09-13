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
		Reward: 99,
		Height: 133853,
		Txs:    []database.DBSimpleTxInBlock{{Fee: 99}},
	}
}

func Test_CreateRetSimpleBlockInfo(t *testing.T) {
	var TxFee int64
	header := newTestDBBlock(t)
	got := createRetSimpleBlockInfo(header)

	for i := 0; i < len(header.Txs); i++ {
		TxFee += header.Txs[i].Fee
	}
	assert.Equal(t, got.Fee, TxFee)
	assert.Equal(t, int64(got.Height), header.Height)
	assert.Equal(t, got.Reward, header.Reward)
}
