package syncer

import (
	"sync"
	"time"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

func (s *Syncer) getAccountFromDBOrCache(address string) *database.DBAccount {
	account, ok := s.cacheAccount[address]
	if ok {
		s.updateAccount[address] = account
		return account
	}

	fromAccount, err := s.db.GetAccountByAddress(address)
	if err != nil {
		fromAccount = database.CreateEmptyAccount(address, s.shardNumber)
	}

	s.updateAccount[address] = fromAccount
	s.cacheAccount[address] = fromAccount
	return fromAccount
}

func (s *Syncer) getMinerAccountAndCount(account *database.DBAccount, reward int64, txFee int64) {
	miner, ok := s.cacheMinerAccount[account.Address]
	if ok {
		miner.Reward += reward
		miner.TxFee += txFee
		miner.Revenue = miner.Reward + miner.TxFee
		s.updateMinerAccount[account.Address] = miner
		return
	}

	minerAccount, err := s.db.GetMinerAccountByAddress(account.Address)
	if err != nil {
		minerAccount = &database.DBMiner{
			Address:     account.Address,
			Revenue:     account.Balance,
			ShardNumber: account.ShardNumber,
			Reward:      reward,
			TxFee:       txFee,
			TimeStamp:   time.Now().Unix(),
		}
	}

	s.cacheMinerAccount[account.Address] = minerAccount
	s.updateMinerAccount[account.Address] = minerAccount
}

//ProcessAccount Process All Account included in the block
func (s *Syncer) accountSync(b *rpc.BlockInfo) error {
	txFees := int64(0)
	for i := 0; i < len(b.Txs); i++ {
		tx := b.Txs[i]
		txFees += tx.Fee
		//exclude coinbase transaction
		if tx.From != nullAddress {

			fromAccount := s.getAccountFromDBOrCache(tx.From)
			fromAccount.TxCount++

			// fromAccount, err := s.db.GetAccountByAddress(tx.From)
			// if err != nil {
			// 	fromAccount = database.CreateEmptyAccount(tx.From, s.shardNumber)
			// 	err := s.db.AddAccount(fromAccount)
			// 	if err != nil {
			// 		log.Error("[DB] err : %v", err)
			// 		continue
			// 	}
			// }

			// balance, err := s.rpc.GetBalance(tx.From)
			// if err != nil {
			// 	log.Error(err)
			// 	balance = 0
			// }

			// txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.From)
			// if err != nil {
			// 	log.Error(err)
			// 	txCnt = 0
			// }

			// s.db.UpdateAccount(tx.From, balance, txCnt)
		}

		if tx.To == "" {
			//create contract transaction
			//Get contract address from receipt

			receipt, err := s.rpc.GetReceiptByTxHash(tx.Hash)
			if err == nil {
				contractAddress := receipt.ContractAddress
				contractAccount := s.getAccountFromDBOrCache(contractAddress)
				//contractAccount := database.CreateEmptyAccount(contractAddress, s.shardNumber)
				contractAccount.AccType = 1
				contractAccount.TxCount++
				// err := s.db.AddAccount(contractAccount)

				// if err != nil {
				// 	log.Error("[DB] err : %v", err)
				// 	continue
				// }

				// balance, err := s.rpc.GetBalance(contractAddress)
				// if err != nil {
				// 	log.Error(err)
				// 	balance = 0
				// }

				// txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, contractAddress)
				// if err != nil {
				// 	log.Error(err)
				// 	txCnt = 0
				// }

				// s.db.UpdateAccount(contractAddress, balance, txCnt)
			}

		} else {
			toAccount := s.getAccountFromDBOrCache(tx.To)
			toAccount.TxCount++
			// toAccount, err := s.db.GetAccountByAddress(tx.To)
			// if err != nil {
			// 	toAccount = database.CreateEmptyAccount(tx.To, s.shardNumber)
			// 	err := s.db.AddAccount(toAccount)
			// 	if err != nil {
			// 		log.Error("[DB] err : %v", err)
			// 		continue
			// 	}
			// }

			// balance, err := s.rpc.GetBalance(tx.To)
			// if err != nil {
			// 	log.Error(err)
			// 	balance = 0
			// }

			// txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.To)
			// if err != nil {
			// 	log.Error(err)
			// 	txCnt = 0
			// }

			// s.db.UpdateAccount(tx.To, balance, txCnt)
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
			}
		}

		blockCnt, err := s.db.GetMinedBlocksCntByShardNumberAndAddress(s.shardNumber, b.Creator)
		if err != nil {
			log.Error(err)
			blockCnt = 0
		}

		s.db.UpdateAccountMinedBlock(b.Creator, blockCnt)

		s.getMinerAccountAndCount(minerAccount, b.Txs[0].Amount.Int64(), txFees)
	}

	return nil
}

func (s *Syncer) accountUpdateSync() {
	// for _, v := range s.newAccount {
	// 	balance, err := s.rpc.GetBalance(v.Address)
	// 	if err != nil {
	// 		log.Error(err)
	// 		balance = 0
	// 	}
	// 	v.Balance = balance

	// 	// txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, v.Address)
	// 	// if err != nil {
	// 	// 	log.Error(err)
	// 	// 	txCnt = 0
	// 	// }
	// 	// v.TxCount = txCnt

	// 	err = s.db.AddAccount(v)
	// 	if err != nil {
	// 		log.Error("[DB] err : %v", err)
	// 		continue
	// 	}
	// }

	var wg sync.WaitGroup
	wg.Add(len(s.updateAccount) + len(s.updateMinerAccount))

	for _, v := range s.updateAccount {

		txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, v.Address)
		if err != nil {
			log.Error(err)
			txCnt = 0
		}

		v.TxCount = txCnt

		v := v

		s.workerpool.Submit(func() {

			account := v
			balance, err := s.rpc.GetBalance(account.Address)
			if err != nil {
				log.Error(err)
				balance = 0
			}

			account.Balance = balance

			s.db.UpdateAccount(account)

			wg.Done()
		})
	}

	for _, m := range s.updateMinerAccount {
		s.workerpool.Submit(func() {
			s.db.UpdateMinerAccount(m)
			wg.Done()
		})
	}

	wg.Wait()

	s.updateAccount = make(map[string]*database.DBAccount)
	s.updateMinerAccount = make(map[string]*database.DBMiner)
}
