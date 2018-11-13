package syncer

import (
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

const (
	maxInsertConn = 200
)

// Syncer is Seele synchronization handler
type Syncer struct {
	rpc         *rpc.SeeleRPC
	db          Database
	shardNumber int
	syncCnt     int
	workerpool  *workerpool.WorkerPool

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

			blockCnt, err := s.db.GetMinedBlocksCntByShardNumberAndAddress(s.shardNumber, dbBlock.Creator)
			if err != nil {
				log.Error(err)
				blockCnt = 0
			}

			s.db.UpdateAccountMinedBlock(dbBlock.Creator, blockCnt)
		}

		fallBack = true
	}
	log.Info("checkOlderBlocks end-------")
	return fallBack
}

// sync get block data from seele node and store it in the mongodb
func (s *Syncer) sync() error {
	log.Info("[BlockSync syncCnt:%d]Begin Sync", s.syncCnt)
	s.checkOlderBlocks()
	// get seele node block height
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

	log.Info("sync begin-------")
	log.Info("sync dbBlockHeight[%d]", dbBlockHeight)
	for i := dbBlockHeight; i <= curHeight; i++ {
		log.Info("begin to sync block[%d]:", i)
		if s.SyncHandle(i) {
			log.Info("failed to sync block[%d]:", i)
			break
		}
		log.Info("successfully to sync block[%d]:", i)
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
	if err = s.blockSync(rpcBlock); err != nil {
		return true
	}

	// sync transactions
	if err = s.txSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}

	// sync debts
	if err = s.debttxSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}

	// sync accounts
	if err = s.accountSync(rpcBlock); err != nil {
		log.Error(err)
		return true
	}
	s.accountUpdateSync()
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
