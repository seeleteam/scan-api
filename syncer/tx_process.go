package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

// txSync insert the transactions into database
func (s *Syncer) txSync(block *rpc.BlockInfo) error {
	transIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)
	txs := []interface{}{}
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

		// transaction fee is in the receipt
		receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
		if err == nil {
			dbTx.Fee = receipt.TotalFee
		}
		dbTx.Fee = receipt.TotalFee
		if trans.To == "" {
			dbTx.TxType = 1
			dbTx.ContractAddress = receipt.ContractAddress
			dbTx.Receipt = *receipt
		}

		txs = append(txs, dbTx)
	}

	if len(txs) == 0 {
		return nil
	}

	return s.db.AddTxs(txs...)
}

// txcountSync insert the Total number of transactions into database
func (s *Syncer) txcountSync(block *rpc.BlockInfo) error {
	s.db.RemoveAllTxsCount()
	txcount, _ := s.db.GetTxCount()
	dbTxCount := database.CreateDbTxCount(txcount)

	return s.db.AddTxsCount(dbTxCount)
}

// debttxSync insert the debt into database
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

// pendingTxsSync insert the pending transactions into database
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
