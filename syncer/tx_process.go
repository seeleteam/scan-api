package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) txSync(block *rpc.BlockInfo) error {
	for j := 0; j < len(block.Txs); j++ {
		trans := block.Txs[j]
		trans.Block = block.Height
		transIdx, err := s.db.GetTxCntByShardNumber(s.shardNumber)

		//must be an create contract transaction
		if trans.To == "" {

			trans.TxType = 1

			receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
			if err != nil {
				trans.To = receipt.ContractAddress
			}
		}

		if err == nil {
			trans.Idx = transIdx
			dbTx := database.CreateDbTx(trans)
			dbTx.Pending = false
			dbTx.ShardNumber = s.shardNumber
			s.db.AddTx(dbTx)
		} else {
			return err
		}
	}

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
