package syncer

import (
	"sync"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) txSync(block *rpc.BlockInfo) error {
	transIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)

	var wg sync.WaitGroup
	wg.Add(len(block.Txs))

	for j := 0; j < len(block.Txs); j++ {
		trans := block.Txs[j]
		trans.Block = block.Height

		transIdx++
		trans.Idx = transIdx
		dbTx := database.CreateDbTx(trans)
		dbTx.Pending = false
		dbTx.ShardNumber = s.shardNumber

		//must be an create contract transaction
		if trans.To == "" {
			dbTx.TxType = 1
			receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
			if err == nil {
				dbTx.ContractAddress = receipt.ContractAddress
				dbTx.Receipt = *receipt
			}
		}

		s.workerpool.Submit(func() {
			s.db.AddTx(dbTx)
			wg.Done()
		})
	}

	wg.Wait()

	return nil
}

func (s *Syncer) pendingTxsSync() error {
	s.db.RemoveAllPendingTxs()
	transIdx, err := s.db.GetPendingTxCntByShardNumber(s.shardNumber)

	txs, err := s.rpc.GetPendingTransactions()
	if err != nil {
		log.Error(err)
		return err
	}

	for i := 0; i < len(txs); i++ {
		transIdx++
		txs[i].Idx = transIdx
		dbTx := database.CreateDbTx(txs[i])
		dbTx.ShardNumber = s.shardNumber
		dbTx.Pending = true
		err := s.db.AddPendingTx(dbTx)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}
