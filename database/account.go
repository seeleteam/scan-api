/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"scan-api/log"
	"scan-api/rpc"
	"sync"
)

const (
	maxShowAccountNum   = 10000
	maxStoreTxInAccount = 25
	nullAddress         = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
)

var (
	gAccountTbl []*DBAccount
	accMutex    sync.RWMutex
)

//createEmptyAccount create an empty dbaccount
func createEmptyAccount(address string) *DBAccount {
	return &DBAccount{
		Address: address,
	}
}

//createDBAccountTx create an dbaccounttx
func createDBAccountTx(b *rpc.BlockInfo, rpcTx *rpc.Transaction, inOrOut bool) *DBAccountTx {
	tx := &DBAccountTx{
		Hash:      rpcTx.Hash,
		Block:     int64(b.Height),
		From:      rpcTx.From,
		To:        rpcTx.To,
		Amount:    rpcTx.Amount.Int64(),
		Timestamp: int64(rpcTx.Timestamp),
		TxFee:     0.0,
		InOrOut:   inOrOut,
	}

	return tx
}

//ProcessAccount Process All Account included in the block
func ProcessAccount(b *rpc.BlockInfo) {
	for i := 0; i < len(b.Txs); i++ {
		tx := b.Txs[i]

		if tx.From != nullAddress {
			fromAccount, err := GetAccountByAddress(tx.From)
			if err != nil {
				fromAccount = createEmptyAccount(tx.From)
				err := AddAccount(fromAccount)
				if err != nil {
					log.Error("[DB] err : %v", err)
					continue
				} else {
					fromAccount, err = GetAccountByAddress(tx.From)
					if err != nil {
						log.Error("[DB] err : %v", err)
						continue
					}
				}
			}

			var txs []DBAccountTx
			dbAccTx := createDBAccountTx(b, &tx, false)
			txs = append(txs, *dbAccTx)
			txlen := len(fromAccount.Txs)
			if txlen >= maxStoreTxInAccount {
				txlen--
			}

			for i := 0; i < txlen; i++ {
				txs = append(txs, fromAccount.Txs[i])
			}
			UpdateAccount(tx.From, float64(-tx.Amount.Int64()), &txs)
		}

		toAccount, err := GetAccountByAddress(tx.To)
		if err != nil {
			toAccount = createEmptyAccount(tx.To)
			err := AddAccount(toAccount)
			if err != nil {
				log.Error("[DB] err : %v", err)
				continue
			} else {
				toAccount, err = GetAccountByAddress(tx.To)
				if err != nil {
					log.Error("[DB] err : %v", err)
					continue
				}
			}
		}

		var txs []DBAccountTx
		dbAccTx := createDBAccountTx(b, &tx, true)
		txs = append(txs, *dbAccTx)
		txlen := len(toAccount.Txs)
		if txlen >= maxStoreTxInAccount {
			txlen--
		}

		for i := 0; i < txlen; i++ {
			txs = append(txs, toAccount.Txs[i])
		}

		UpdateAccount(tx.To, float64(tx.Amount.Int64()), &txs)
	}

	//exclude genesis block
	if b.Creator != nullAddress {
		minerAccount, err := GetAccountByAddress(b.Creator)
		if err != nil {
			minerAccount = createEmptyAccount(b.Creator)
			err := AddAccount(minerAccount)
			if err != nil {
				log.Error("[DB] err : %v", err)
			} else {
				minerAccount, err = GetAccountByAddress(b.Creator)
				if err != nil {
					log.Error("[DB] err : %v", err)
				}
			}
		}
		UpdateAccountMinedBlock(b.Creator, 1)
	}
}

//ProcessGAccountTable process global account table
func ProcessGAccountTable() {
	temp, err := GetAccounts(maxShowAccountNum)
	if err != nil {
		log.Error("[DB] err : %v", err)
	} else {
		accMutex.Lock()
		gAccountTbl = temp
		accMutex.Unlock()
	}
}

//GetAccountCnt get the length of account table
func GetAccountCnt() int {
	accMutex.RLock()
	size := len(gAccountTbl)
	accMutex.RUnlock()
	return size
}

//GetAccountsByIdx get a transaction list from mongo by time period
func GetAccountsByIdx(begin uint64, end uint64) []*DBAccount {
	accMutex.RLock()
	if end > uint64(len(gAccountTbl)) {
		end = uint64(len(gAccountTbl)) - 1
	}

	retAccounts := gAccountTbl[begin:end]
	accMutex.RUnlock()
	return retAccounts
}
