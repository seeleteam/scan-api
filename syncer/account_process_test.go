/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package syncer

import (
	"testing"

	"github.com/seeleteam/scan-api/database"
	"github.com/stretchr/testify/assert"
)

func newTestDBAccount(t *testing.T) *database.DBAccount {
	return &database.DBAccount{
		ShardNumber: 1,
		Address:     "0x4d7adfbeabc99303c9ad4fcd6370b611b88eb5b3b69b1dbc61f6c0ce30a352b0",
	}
}

func TestCreateDbBlock(t *testing.T) {
	var reward, txFee int64
	s := &Syncer{}
	AccountInfo := newTestDBAccount(t)
	s.cacheMinerAccount = make(map[string]*database.DBMiner)
	s.updateMinerAccount = make(map[string]*database.DBMiner)

	minerAccount := &database.DBMiner{
		Address:     AccountInfo.Address,
		Revenue:     AccountInfo.Balance,
		ShardNumber: AccountInfo.ShardNumber,
		Reward:      reward,
		TxFee:       txFee,
	}

	s.cacheMinerAccount[AccountInfo.Address] = minerAccount
	s.getMinerAccountAndCount(AccountInfo, reward, txFee)
	assert.Equal(t, AccountInfo.Address, minerAccount.Address)
	assert.Equal(t, AccountInfo.ShardNumber, minerAccount.ShardNumber)
}
