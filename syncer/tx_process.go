package syncer

import (
	"strconv"
	"strings"
	"time"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

// txSync insert the transactions into database
func (s *Syncer) txSync(block *rpc.BlockInfo) error {
	transIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)
	txs := []interface{}{}
	var dbTxs []*database.DBTx
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
			if trans.To == "" {
				dbTx.TxType = 1
				dbTx.ContractAddress = receipt.ContractAddress
				dbTx.Receipt = *receipt
			}
		}
		dbTxs = append(dbTxs, dbTx)
		txs = append(txs, dbTx)
	}

	if len(txs) == 0 {
		return nil
	}

	if err := s.db.AddTxs(txs...); err != nil {
		return err
	}

	return nil
}

func dropOutDate(db Database, today, startDate string) {
	if hasOutDate(db, today) {
		db.RemoveOutDateByDate(startDate)
	}
}

func hasOutDate(db Database, today string) bool {
	cnt, err := db.GetTxHisCntByDate(today)
	if err != nil {
		return false
	}
	if cnt == 0 {
		return true
	}
	return false
}

func nextDate(date *string) {
	ymd := strings.Split(*date, "-")
	year, _ := strconv.Atoi(ymd[0])
	month, _ := strconv.Atoi(ymd[1])
	day, _ := strconv.Atoi(ymd[2])
	dateTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dateTime = dateTime.AddDate(0, 0, 1)
	*date = dateTime.Format("2006-01-02")
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
	transIdx, _ := s.db.GetPendingTxCntByShardNumber(s.shardNumber)

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
