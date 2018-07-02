/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package txhistory

import (
	"sync"
	"time"

	"github.com/seeleteam/scan-api/chart"
	"github.com/seeleteam/scan-api/database"

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
		for i := 1; i <= chart.ShardCount; i++ {
			ProcessOneDayTransaction(i, next)
		}
	}
}

//ProcessOldTransactions Count transactions happened in the past days
func ProcessOldTransactions() {
	for i := 1; i <= chart.ShardCount; i++ {
		now := time.Now()
		//lastZeroTime := now.Add(-time.Hour * 24)
		todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		for {

			lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
			_, err := chart.GChartDB.GetOneDayTransInfo(i, lastZeroTime.Unix())
			if err != mgo.ErrNotFound {
				break
			}
			if !ProcessOneDayTransaction(i, todayZeroTime) {
				break
			}
			todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
		}
	}
}

//ProcessOneDayTransaction Count transactions that occur within one day
func ProcessOneDayTransaction(shardNumber int, day time.Time) bool {
	secondBlock, err := chart.GChartDB.GetBlockByHeight(shardNumber, 1)
	if err != nil {
		return false
	}

	if day.Unix() < secondBlock.Timestamp {
		return false
	}

	thisZeroTime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	lastDayZeroTime := day.Add(-time.Hour * 24)
	var dbBlocks []*database.DBBlock
	dbBlocks, err = chart.GChartDB.GetBlocksByTime(shardNumber, lastDayZeroTime.Unix(), thisZeroTime.Unix())
	if err != nil {
		return false
	}

	var info database.DBOneDayTxInfo
	for i := 0; i < len(dbBlocks); i++ {
		info.TotalBlocks++
		info.TotalTxs += len(dbBlocks[i].Txs)
	}
	info.TimeStamp = lastDayZeroTime.Unix()
	chart.GChartDB.AddOneDayTransInfo(shardNumber, &info)
	return true
}

func init() {
	chart.RegisterProcessFunc(Process)
}
