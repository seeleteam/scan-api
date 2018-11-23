package syncer

import (
	"fmt"
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
	s.mu.Lock()
	defer s.mu.Unlock()
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
	//s.mu.Lock()
	// defer s.mu.Unlock()
	miner, ok := s.cacheMinerAccount[account.Address]
	if ok {
		miner.Reward += reward
		miner.TxFee += txFee
		miner.Revenue = miner.Reward + miner.TxFee
		s.updateMinerAccount[account.Address] = miner
		return
	}

	s.getMinerAccount(account, reward, txFee)
}

func (s *Syncer) getMinerAccount(account *database.DBAccount, reward int64, txFee int64) {
	s.mu.Lock()

	minerAccount, err := s.db.GetMinerAccountByAddress(account.Address)
	if err != nil {
		fmt.Println("===============================================")
		fmt.Println("miner00000000000000")
		fmt.Println("===============================================")
		minerAccount = &database.DBMiner{
			Address:     account.Address,
			ShardNumber: account.ShardNumber,
			TimeStamp:   time.Now().Unix(),
		}
	}

	minerAccount.Reward += reward
	minerAccount.TxFee += txFee
	minerAccount.Revenue = minerAccount.Reward + minerAccount.TxFee
	fmt.Println("===============================================")
	fmt.Println("miner1111111111111111:", minerAccount)
	fmt.Println("===============================================")

	fmt.Println("===============================================")
	fmt.Println("miner22222222222222:", minerAccount)
	fmt.Println("===============================================")
	s.cacheMinerAccount[account.Address] = minerAccount
	s.updateMinerAccount[account.Address] = minerAccount
	defer s.mu.Unlock()

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
		}

		if tx.To == "" {
			//create contract transaction
			//Get contract address from receipt
			receipt, err := s.rpc.GetReceiptByTxHash(tx.Hash)
			if err == nil {
				contractAddress := receipt.ContractAddress
				contractAccount := s.getAccountFromDBOrCache(contractAddress)
				contractAccount.AccType = 1
				contractAccount.TxCount++
			}
		} else {
			toAccount := s.getAccountFromDBOrCache(tx.To)
			toAccount.TxCount++
		}
	}

	//exclude genesis block
	if b.Creator != nullAddress {
		minerAccount := s.getAccountFromDBOrCache(b.Creator)
		blockCnt, err := s.db.GetMinedBlocksCntByShardNumberAndAddress(s.shardNumber, b.Creator)
		if err != nil {
			log.Error(err)
			blockCnt = 0
		}

		minerAccount.Mined = blockCnt
		s.getMinerAccountAndCount(minerAccount, b.Txs[0].Amount.Int64(), txFees)
	}
	return nil
}

func (s *Syncer) accountUpdateSync() {
	fmt.Println("===============================================")
	fmt.Println("s.updateMinerAccount[string(i)]:", s.updateMinerAccount)
	fmt.Println("s.updateAccount[string(i)]:", s.updateAccount)
	fmt.Println("===============================================")
	var mutex sync.Mutex
	mutex.Lock()
	channels := make([]chan int, len(s.updateMinerAccount))
	for i := 0; i < len(s.updateMinerAccount); i++ {
		fmt.Println("s.updateMinerAccount[string(i)]:", s.updateMinerAccount[string(i)])
		channels[i] = make(chan int)
		go func(i int, c chan int) {
			mutex.Lock()
			// for _, m := range s.updateMinerAccount {
			s.workerpool.Submit(func() {
				s.db.UpdateMinerAccount(s.updateMinerAccount[string(i)])
			})

			// end

			mutex.Unlock()
			c <- i
		}(i, channels[i])
	}
	mutex.Unlock()
	for _, c := range channels {
		<-c
	}
	s.updateAccount = make(map[string]*database.DBAccount)
	//	s.updateMinerAccount = make(map[string]*database.DBMiner)
}

// func (s *Syncer) accountUpdateSync() {
// 	var wg sync.WaitGroup
// 	wg.Add(len(s.updateAccount) + len(s.updateMinerAccount))

// 	for _, v := range s.updateAccount {

// 		txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, v.Address)
// 		if err != nil {
// 			log.Error(err)
// 			txCnt = 0
// 		}

// 		v.TxCount = txCnt

// 		v := v

// 		s.workerpool.Submit(func() {

// 			account := v
// 			balance, err := s.rpc.GetBalance(account.Address)
// 			if err != nil {
// 				log.Error(err)
// 				balance = 0
// 			}

// 			account.Balance = balance

// 			s.db.UpdateAccount(account)

// 			wg.Done()
// 		})
// 	}

// 	for _, m := range s.updateMinerAccount {
// 		s.workerpool.Submit(func() {
// 			s.db.UpdateMinerAccount(m)
// 			wg.Done()
// 		})
// 	}

// 	wg.Wait()

// 	s.updateAccount = make(map[string]*database.DBAccount)
// 	s.updateMinerAccount = make(map[string]*database.DBMiner)
// }
