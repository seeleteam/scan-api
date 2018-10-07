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
	s.db.AddTxs(txs...)
	return nil
}

func (s *Syncer) debttxSync(block *rpc.BlockInfo) error {
	debttxs := []interface{}{}
	for i := 0; i < len(block.Debts); i++ {
		debts := block.Debts[i]

		for j := 0; j < len(block.TxDebts); j++ {
			if block.Debts[i].TxHash == block.TxDebts[j].TxHash {
				debts.TxHash = block.TxDebts[j].TxHash
			}
		}

		debts.Block = block.Height
		debtTx := database.CreateDebtTx(debts)
		debtTx.ShardNumber = s.shardNumber

		debttxs = append(debttxs, debtTx)

	}
	//wg.Wait()
	s.db.AddDebtTxs(debttxs...)
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
