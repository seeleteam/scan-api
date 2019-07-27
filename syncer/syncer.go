package syncer

import (
	"fmt"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

const (
	maxInsertConn = 20 //200
)

// Syncer is Seele synchronization handler
type Syncer struct {
	rpc                *rpc.SeeleRPC
	db                 Database
	shardNumber        int
	syncCnt            int
	workerpool         *workerpool.WorkerPool
	mu                 sync.Mutex
	cacheAccount       map[string]*database.DBAccount
	updateAccount      map[string]*database.DBAccount
	cacheMinerAccount  map[string]*database.DBMiner
	updateMinerAccount map[string]*database.DBMiner
}

// NewSyncer return a syncer to sync block data from seele node
func NewSyncer(db Database, rpcConnURL string, shardNumber int) *Syncer {
	rpc := rpc.NewRPC(rpcConnURL)
	if rpc == nil {
		return nil
	}

	if err := rpc.Connect(); err != nil {
		fmt.Printf("rpc init failed, connurl:%v\n", rpcConnURL)
		return nil
	}

	return &Syncer{
		db:                 db,
		rpc:                rpc,
		shardNumber:        shardNumber,
		syncCnt:            0,
		cacheAccount:       make(map[string]*database.DBAccount),
		updateAccount:      make(map[string]*database.DBAccount),
		cacheMinerAccount:  make(map[string]*database.DBMiner),
		updateMinerAccount: make(map[string]*database.DBMiner),
		workerpool:         workerpool.New(maxInsertConn),
	}
}

// Blocks that are already in storage may be modified
func (s *Syncer) checkOlderBlocks() bool {
	dbBlockHeight, err := s.db.GetBlockHeight(s.shardNumber)
	if err != nil {
		log.Error(err)
		return false
	}

	if dbBlockHeight == 0 {
		return false
	}

	fallBack := false
	log.Debug("checkOlderBlocks begin-------")
	for i := dbBlockHeight - 1; i >= 0; i-- {
		rpcBlock, err := s.rpc.GetBlockByHeight(i, true)
		if err != nil {
			return fallBack
		}

		dbBlock, err := s.db.GetBlockByHeight(s.shardNumber, i)
		if err != nil {
			return fallBack
		}

		//if the block data is different wo should sync the data again
		if dbBlock.HeadHash == rpcBlock.Hash {
			return fallBack
		}

		log.Info("checkOlderBlocks remove block [%d]: %v", i, dbBlock)
		//Delete dbBlock
		s.db.RemoveBlock(s.shardNumber, i)

		//Delete txs
		s.db.RemoveTxs(s.shardNumber, i)

		//Modify accounts
		for j := 0; j < len(dbBlock.Txs); j++ {
			tx := dbBlock.Txs[j]

			if tx.To != "" {
				toAccount, err := s.db.GetAccountByAddress(tx.To)
				if err != nil {
					log.Error(err)
					return fallBack
				}

				toAccount.Balance, err = s.rpc.GetBalance(tx.To)
				if err != nil {
					log.Error(err)
					toAccount.Balance = 0
				}

				txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.To)
				if err != nil {
					log.Error(err)
					txCnt = 0
				}

				toAccount.TxCount = int64(txCnt)
				s.db.UpdateAccount(toAccount)
			} else {
				receipt, err := s.rpc.GetReceiptByTxHash(tx.Hash)
				if err == nil {
					contractAddress := receipt.ContractAddress

					contractAccount, err := s.db.GetAccountByAddress(contractAddress)
					if err != nil {
						log.Error(err)
						return fallBack
					}

					contractAccount.Balance, err = s.rpc.GetBalance(contractAddress)
					if err != nil {
						log.Error(err)
						contractAccount.Balance = 0
					}

					txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, contractAddress)
					if err != nil {
						log.Error(err)
						txCnt = 0
					}

					contractAccount.TxCount = int64(txCnt)
					s.db.UpdateAccount(contractAccount)
				}
			}

			if tx.From != nullAddress {
				fromAccount, err := s.db.GetAccountByAddress(tx.From)

				fromAccount.Balance, err = s.rpc.GetBalance(tx.From)
				if err != nil {
					log.Error(err)
					fromAccount.Balance = 0
				}

				txCnt, err := s.db.GetTxCntByShardNumberAndAddress(s.shardNumber, tx.From)
				if err != nil {
					log.Error(err)
					txCnt = 0
				}

				fromAccount.TxCount = int64(txCnt)
				s.db.UpdateAccount(fromAccount)
			}
		}

		if dbBlock.Creator != nullAddress {
			minerAccount, err := s.db.GetAccountByAddress(dbBlock.Creator)
			if err != nil {
				minerAccount = database.CreateEmptyAccount(dbBlock.Creator, s.shardNumber)
				err := s.db.AddAccount(minerAccount)
				if err != nil {
					log.Error("[DB] err : %v", err)
				} else {
					minerAccount, err = s.db.GetAccountByAddress(dbBlock.Creator)
					if err != nil {
						log.Error("[DB] err : %v", err)
					}
				}
			}

		}

		fallBack = true
	}
	log.Info("checkOlderBlocks end-------")
	return fallBack
}

var wg sync.WaitGroup

// sync get block data from seele node and store it in the mongodb
func (s *Syncer) sync() error {
	log.Info("[BlockSync syncCnt:%d]Begin Sync", s.syncCnt)
	s.checkOlderBlocks()
	// get seele node block height
ErrContinue:
	curHeight, err := s.rpc.CurrentBlockHeight()
	if err != nil {
		log.Error(err)
		return err
	}

	// get local block height
	dbBlockHeight, err := s.db.GetBlockHeight(s.shardNumber)
	if err != nil {
		log.Error(err)
		return err
	}
	if curHeight <= dbBlockHeight || (curHeight - dbBlockHeight) < 100 {
		log.Info("not enough block to sync")
		return nil
	}
	log.Info("sync begin-------")
	log.Info("sync dbBlockHeight[%d]", dbBlockHeight)
	maxSyncCnt :=(uint64(1)) // use single thread; 200
	anum := curHeight - dbBlockHeight
	if anum <= 0 {
		log.Info("block chain height is smaller than scan db block height")
		return nil
	}
	if anum >= maxSyncCnt {
		anum = maxSyncCnt
	}
	wg.Add(int(anum))
	abc := dbBlockHeight + anum
	var i uint64
	for i = dbBlockHeight; i < abc; i++ {
		log.Info("begin to sync block[%d]:", i)
		go func(i uint64) {
			defer wg.Done()
			result := s.SyncHandle(i)
			if(result){
				log.Error("sync block [%d] failed",i)
			}else{
				log.Info("successfully to sync block[%d]:", i)
			}
		}(i)
	}
	wg.Wait()
	if anum >= maxSyncCnt {
		goto ErrContinue
	}
	log.Info("sync end-------")

	err = s.pendingTxsSync()
	if err != nil {
		log.Error(err)
	}
	log.Info("[BlockSync syncCnt:%d]End Sync", s.syncCnt)
	s.syncCnt++
	return nil
}

// SyncHandle sync the block data from seele node, and handle tx or account
func (s *Syncer) SyncHandle(i uint64) bool {
	rpcBlock, err := s.rpc.GetBlockByHeight(i, true)
	if err != nil {
		s.rpc.Release()
		log.Error(err)
		return true
	}

	// sync block
	timeBegin := time.Now().Unix()
	if err = s.blockSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	log.Debug("syncerHandle blockSync time: %d(s)",time.Now().Unix()-timeBegin)

	// sync transactions
	timeBegin = time.Now().Unix()
	if err = s.txSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	log.Debug("syncerHandle txSync time: %d(s)",time.Now().Unix()-timeBegin)

	// sync debts
	timeBegin = time.Now().Unix()
	if err = s.debttxSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	log.Debug("syncerHandle debttxSync time: %d(s)",time.Now().Unix()-timeBegin)
	// sync accounts
	timeBegin = time.Now().Unix()
	if err = s.accountSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	log.Debug("syncerHandle accountSync time: %d(s)",time.Now().Unix()-timeBegin)
	// sync minersaccount
	timeBegin = time.Now().Unix()
	if err = s.minersaccountSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	log.Debug("syncerHandle minersaccountSync time: %d(s)",time.Now().Unix()-timeBegin)
	return false
}

// StartSync start an timer to sync block data from seele node
func (s *Syncer) StartSync(interval time.Duration) {
	s.sync()

	ticks := time.NewTicker(interval * time.Second)
	tick := ticks.C
	i := 0
	go func() {
		for range tick {
			log.Info("StartSync[%d].............", i)
			s.sync()
			i++
			_, ok := <-tick
			if !ok {
				break
			}
		}
	}()
}
