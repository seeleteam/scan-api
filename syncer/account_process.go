package syncer

import (
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"

	"github.com/seeleteam/scan-api/database"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

//ProcessAccount Process All Account included in the block
func (s *Syncer) accountSync(b *rpc.BlockInfo) error {
	for i := 0; i < len(b.Txs); i++ {
		tx := b.Txs[i]

		//exclude coinbase transaction
		if tx.From != nullAddress {
			fromAccount, err := s.db.GetAccountByAddress(tx.From)
			if err != nil {
				fromAccount = database.CreateEmptyAccount(tx.From, s.shardNumber)
				err := s.db.AddAccount(fromAccount)
				if err != nil {
					log.Error("[DB] err : %v", err)
					continue
				} else {
					fromAccount, err = s.db.GetAccountByAddress(tx.From)
					if err != nil {
						log.Error("[DB] err : %v", err)
						continue
					}
				}
			}

			balance, err := s.rpc.GetBalance(tx.From)
			if err != nil {
				log.Error(err)
				balance = 0
			}

			txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.From)
			if err != nil {
				log.Error(err)
				txCnt = 0
			}

			s.db.UpdateAccount(tx.From, balance, txCnt)
		}

		if tx.To == nullAddress {
			//create contract transaction
			//Get contract address from receipt

			receipt, err := s.rpc.GetReceiptByTxHash(tx.Hash)
			if err != nil {
				contractAddress := receipt.ContractAddress
				contractAccount := database.CreateEmptyAccount(contractAddress, s.shardNumber)
				contractAccount.AccType = 1
				err := s.db.AddAccount(contractAccount)
				if err != nil {
					log.Error("[DB] err : %v", err)
					continue
				} else {
					contractAccount, err = s.db.GetAccountByAddress(contractAddress)
					if err != nil {
						log.Error("[DB] err : %v", err)
						continue
					}
				}

				balance, err := s.rpc.GetBalance(contractAddress)
				if err != nil {
					log.Error(err)
					balance = 0
				}

				txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.From)
				if err != nil {
					log.Error(err)
					txCnt = 0
				}

				s.db.UpdateAccount(contractAddress, balance, txCnt)
			}

		} else {
			toAccount, err := s.db.GetAccountByAddress(tx.To)
			if err != nil {
				toAccount = database.CreateEmptyAccount(tx.To, s.shardNumber)
				err := s.db.AddAccount(toAccount)
				if err != nil {
					log.Error("[DB] err : %v", err)
					continue
				} else {
					toAccount, err = s.db.GetAccountByAddress(tx.To)
					if err != nil {
						log.Error("[DB] err : %v", err)
						continue
					}
				}
			}

			balance, err := s.rpc.GetBalance(tx.To)
			if err != nil {
				log.Error(err)
				balance = 0
			}

			txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.From)
			if err != nil {
				log.Error(err)
				txCnt = 0
			}

			s.db.UpdateAccount(tx.To, balance, txCnt)
		}
	}

	//exclude genesis block
	if b.Creator != nullAddress {
		minerAccount, err := s.db.GetAccountByAddress(b.Creator)
		if err != nil {
			minerAccount = database.CreateEmptyAccount(b.Creator, s.shardNumber)
			err := s.db.AddAccount(minerAccount)
			if err != nil {
				log.Error("[DB] err : %v", err)
			} else {
				minerAccount, err = s.db.GetAccountByAddress(b.Creator)
				if err != nil {
					log.Error("[DB] err : %v", err)
				}
			}
		}

		blockCnt, err := s.db.GetMinedBlocksCntByShardNumberAndAddress(s.shardNumber, b.Creator)
		if err != nil {
			log.Error(err)
			blockCnt = 0
		}

		s.db.UpdateAccountMinedBlock(b.Creator, blockCnt)
	}

	return nil
}
