package syncer

import (
	"sync"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

func (s *Syncer) accountUpdateSync() {
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

func (s *Syncer) accountSync(b *rpc.BlockInfo) error {
	var address string
	var AccType int
	for i := 0; i < len(b.Txs); i++ {
		s.mu.Lock()
		tx := b.Txs[i]
		if tx.From != nullAddress {
			address = tx.From
		}

		if tx.To == "" {
			//create contract transaction
			//Get contract address from receipt
			receipt, err := s.rpc.GetReceiptByTxHash(tx.Hash)
			if err == nil {
				contractAddress := receipt.ContractAddress
				address = contractAddress
				AccType = 1
			}
		} else {
			address = tx.To
		}

		balance, err := s.rpc.GetBalance(address)
		if err != nil {
			log.Error(err)
			balance = 0
		}

		txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, address)
		if err != nil {
			log.Error(err)
			txCnt = 0
		}
		accounts := &database.DBAccount{
			AccType:     AccType,
			ShardNumber: s.shardNumber,
			Address:     address,
			TxCount:     txCnt,
			Balance:     balance,
			TimeStamp:   b.Timestamp.Int64(),
		}

		s.db.UpdateAccount(accounts)
		defer s.mu.Unlock()
	}
	return nil

}

func (s *Syncer) minersaccountSync(b *rpc.BlockInfo) error {
	//exclude genesis block
	if b.Creator != nullAddress {
		s.mu.Lock()
		blockCnt, blockFee, blockAmount, err := s.db.GetMinedBlocksByShardNumberAndAddress(s.shardNumber, b.Creator)
		if err != nil {
			log.Error(err)
			blockCnt = 0
		}

		miners := &database.DBMiner{
			ShardNumber: s.shardNumber,
			Address:     b.Creator,
			Reward:      blockAmount,
			TxFee:       blockFee,
			Revenue:     blockAmount + blockFee,
			TimeStamp:   b.Timestamp.Int64(),
			Mined:       blockCnt,
		}
		s.db.UpdateMinerAccount(miners)
		defer s.mu.Unlock()
	}
	return nil

}
