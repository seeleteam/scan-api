package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) txSync(block *rpc.BlockInfo) error {
	transIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)
	txs := []interface{}{}
	//var wg sync.WaitGroup
	//wg.Add(len(block.Txs))

	for i := 0; i < len(block.Txs); i++ {
		trans := block.Txs[i]
		for j := 0; j < len(block.TxDebts); j++ {
			if block.Txs[i].Hash == block.TxDebts[j].TxHash {
				trans.DebtTxHash = block.TxDebts[j].Hash
			}
		}

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

		txs = append(txs, dbTx)
		//s.workerpool.Submit(func() {
		//	s.db.AddTx(dbTx)
		//	wg.Done()
		//})
	}

	//wg.Wait()
	if len(txs) == 0 {
		return nil
	}

	return s.db.AddTxs(txs...)
}

func (s *Syncer) debttxSync(block *rpc.BlockInfo) error {
	debtIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)
	debttxs := []interface{}{}
	for i := 0; i < len(block.Debts); i++ {
		debts := block.Debts[i]
		debts.Block = block.Height
		debtIdx++
		debts.Idx = debtIdx
		debtTx := database.CreateDebtTx(debts)
		debtTx.ShardNumber = s.shardNumber

		debttxs = append(debttxs, debtTx)
	}

	if len(debttxs) == 0 {
		return nil
	}

	return s.db.AddDebtTxs(debttxs...)
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
