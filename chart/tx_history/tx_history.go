/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package txhistory

import (
	"scan-api/database"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

//Process set an timer to count transactions every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldTransactions()
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		ProcessOneDayTransaction(next)
	}
}

//ProcessOldTransactions Count transactions happened in the past days
func ProcessOldTransactions() {
	now := time.Now()
	//lastZeroTime := now.Add(-time.Hour * 24)
	todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
		_, err := database.GetOneDayTransInfo(lastZeroTime.Unix())
		if err != mgo.ErrNotFound {
			break
		}
		if !ProcessOneDayTransaction(todayZeroTime) {
			break
		}
		todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
	}
}

//ProcessOneDayTransaction Count transactions that occur within one day
func ProcessOneDayTransaction(day time.Time) bool {
	secondBlock, err := database.GetBlockByHeight(1)
	if err != nil {
		return false
	}

	if day.Unix() < secondBlock.Timestamp {
		return false
	}

	thisZeroTime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	lastDayZeroTime := day.Add(-time.Hour * 24)
	var dbBlocks []*database.DBBlock
	dbBlocks, err = database.GetBlocksByTime(lastDayZeroTime.Unix(), thisZeroTime.Unix())
	if err != nil {
		return false
	}

	var info database.DBOneDayTxInfo
	for i := 0; i < len(dbBlocks); i++ {
		info.TotalBlocks++
		info.TotalTxs += len(dbBlocks[i].Txs)
	}
	info.TimeStamp = lastDayZeroTime.Unix()
	database.AddOneDayTransInfo(&info)
	return true
}
