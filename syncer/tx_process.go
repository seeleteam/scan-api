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
	timeBegin := time.Now().Unix()
	transIdx, _ := s.db.GetTxCntByShardNumber(s.shardNumber)
	log.Debug("seele_syncer tx_process GetTxCntByShardNumber time:%d(s)",time.Now().Unix()-timeBegin )
	txs := []interface{}{}
	var dbTxs []*database.DBTx
	timeBegin = time.Now().Unix()
	for i := 0; i < len(block.Txs); i++ {
		trans := block.Txs[i]
		trans.Timestamp = block.Timestamp.Uint64()
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
			dbTx.UsedGas = receipt.UsedGas
			if trans.To == "" {
				dbTx.TxType = 1
				dbTx.ContractAddress = receipt.ContractAddress
				dbTx.Receipt = *receipt
			}
		}
		dbTxs = append(dbTxs, dbTx)
		txs = append(txs, dbTx)
	}
	log.Debug("seele_syncer tx_process prepareTxs time:%d(s)",time.Now().Unix()-timeBegin )

	if len(txs) == 0 {
		return nil
	}
	timeBegin = time.Now().Unix()
	if err := s.db.AddTxs(txs...); err != nil {
		return err
	}
	log.Debug("seele_syncer tx_process AddTxs time:%d(s)",time.Now().Unix()-timeBegin )

	// insert 30 days history transaction number into database
	//s.txHisSync(dbTxs)    // this long-time process should be done in another way

	return nil
}

func (s *Syncer) txHisSync(txs []*database.DBTx) error {
	now := time.Now()
	// get the start date of 30 days history
	startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
	// get which day the txs belong to
	dates := getTxDate(txs)
	// filter the dates, farther than startDate will be filtered out
	dates = filterDate(dates, startDate)
	todayDate := now.Format("2006-01-02")
	// update the transactions count of the days
	s.mu.Lock()
	defer s.mu.Unlock()
	updateTxHisForDates(s.db, dates, todayDate, startDate)
	// if the history number is not 30, insert the other days after start day
	checkTxHis(s.db, startDate, todayDate)
	return nil
}

func getTxDate(txs []*database.DBTx) []string {
	dateTxs := make(map[string]bool)
	for _, tx := range txs {
		dateTxs[tx.Timetxs] = true
	}
	dates := make([]string, 0, len(dateTxs))
	for date := range dateTxs {
		dates = append(dates, date)
	}
	return dates
}

func filterDate(dates []string, limit string) []string {
	var validDates []string
	for _, date := range dates {
		if date >= limit {
			validDates = append(validDates, date)
		}
	}
	return validDates
}

func updateTxHisForDates(db Database, dates []string, today, startDate string) {
	for _, date := range dates {
		if date == today {
			dropOutDate(db, today, startDate)
		}
		updateTxHis(db, date)
	}
}

func updateTxHis(db Database, date string) {
	tx := new(database.DBSimpleTxs)
	var cnt, sumgas int64
	tx.Stime = date
	highGasPrice, lowGasPrice, cnt, sumgas, err := db.GetTxsinfoByDate(date)
	if err != nil {
		cnt = 0
	}
	tx.TxCount = int(cnt)
	tx.GasPrice = sumgas
	tx.HighGasPrice = highGasPrice
	tx.LowGasPrice = lowGasPrice
	db.UpdateTxsCntByDate(tx)
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

func checkTxHis(db Database, startDate, today string) {
	txs, err := db.GetTxHis(startDate, today)
	if err != nil {
		return
	}
	if len(txs) == 30 {
		return
	}
	for date := startDate; date < today; nextDate(&date) {
		updateTxHis(db, date)
	}
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
