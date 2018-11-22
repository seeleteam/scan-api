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
	var wga sync.WaitGroup
	wga.Add(1)
	account, ok := s.cacheAccount[address]
	if ok {
		s.updateAccount[address] = account
		return account
	}

	fromAccount, err := s.db.GetAccountByAddress(address)
	if err != nil {
		fromAccount = database.CreateEmptyAccount(address, s.shardNumber)
	}
	wga.Done()
	s.updateAccount[address] = fromAccount
	s.cacheAccount[address] = fromAccount
	wga.Wait()
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

	s.getMinerAccount(account, reward, txFee)
}

func (s *Syncer) getMinerAccount(account *database.DBAccount, reward int64, txFee int64) {
	minerAccount, err := s.db.GetMinerAccountByAddress(account.Address)
	if err != nil {
		minerAccount = &database.DBMiner{
			Address:     account.Address,
			ShardNumber: account.ShardNumber,
			TimeStamp:   time.Now().Unix(),
		}
	}

	minerAccount.Reward += reward
	minerAccount.TxFee += txFee
	minerAccount.Revenue = minerAccount.Reward + minerAccount.TxFee
	s.cacheMinerAccount[account.Address] = minerAccount
	s.updateMinerAccount[account.Address] = minerAccount
}

//ProcessAccount Process All Account included in the block
func (s *Syncer) accountSync(b *rpc.BlockInfo) error {
	fmt.Println("----bbbbbbbbbb----", b.Height)
	txFees := int64(0)
	var mutex sync.Mutex
	mutex.Lock()
	channels := make([]chan int, len(b.Txs))
	for i := 0; i < len(b.Txs); i++ {
		channels[i] = make(chan int)
		go func(i int, c chan int) {
			mutex.Lock()
			fmt.Println("Locked: ", i)
			//-----------------------------------------------
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
			//-----------------------------------------------
			fmt.Println("Unlock the lock: ", i)
			mutex.Unlock()
			c <- i
		}(i, channels[i])
	}
	//time.Sleep(time.Second)
	fmt.Println("Unlock the lock")
	mutex.Unlock()
	//time.Sleep(time.Second)

	for _, c := range channels {
		fmt.Println("channelschannelschannelschannels: ", channels)
		<-c
	}

	//exclude genesis block
	var wgg sync.WaitGroup
	fmt.Println("[[[[[[[[[[[[[[[[[[wgg]]]]]]]]]]]]]]]]]]]")
	wgg.Add(1)
	fmt.Println("[[[[[[[[[[[[[[[[[[wgg1111111]]]]]]]]]]]]]]]]]]]")
	if b.Creator != nullAddress {
		minerAccount := s.getAccountFromDBOrCache(b.Creator)
		blockCnt, err := s.db.GetMinedBlocksCntByShardNumberAndAddress(s.shardNumber, b.Creator)
		if err != nil {
			log.Error(err)
			blockCnt = 0
		}

		minerAccount.Mined = blockCnt
		s.getMinerAccountAndCount(minerAccount, b.Txs[0].Amount.Int64(), txFees)
		wgg.Done()
	} else {
		wgg.Done()
	}
	wgg.Wait()
	return nil
}

func (s *Syncer) accountUpdateSync() {
	var wg sync.WaitGroup
	wg.Add(len(s.updateAccount) + len(s.updateMinerAccount))
	fmt.Println("---s.updateAccount--------s.updateMinerAccount-------", s.updateAccount, s.updateMinerAccount)
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
