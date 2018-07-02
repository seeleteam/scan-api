/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package block

import (
	"sync"
	"time"

	"github.com/seeleteam/scan-api/chart"
	"github.com/seeleteam/scan-api/database"
	mgo "gopkg.in/mgo.v2"
)

//Process set an timer to process block data every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldBlocks()
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		for i := 1; i <= chart.ShardCount; i++ {
			ProcessOneDayBlocks(i, next)
		}
	}
}

//ProcessOldBlocks Process block data mined in the past
func ProcessOldBlocks() {

	for i := 1; i <= chart.ShardCount; i++ {
		now := time.Now()
		//lastZeroTime := now.Add(-time.Hour * 24)
		todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		for {

			lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
			_, err := chart.GChartDB.GetOneDayBlock(i, lastZeroTime.Unix())
			if err != mgo.ErrNotFound {
				break
			}
			if !ProcessOneDayBlocks(i, todayZeroTime) {
				break
			}
			todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
		}
	}

}

//ProcessOneDayBlocks Process an single day block data
func ProcessOneDayBlocks(shardNumber int, day time.Time) bool {
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

	var info database.DBOneDayBlockInfo

	for i := 0; i < len(dbBlocks); i++ {
		info.TotalBlocks++
		txLen := len(dbBlocks[i].Txs)
		if txLen > 0 {
			tx := dbBlocks[i].Txs[txLen-1]
			info.Rewards += tx.Amount
		}
	}

	info.TimeStamp = lastDayZeroTime.Unix()
	chart.GChartDB.AddOneDayBlock(shardNumber, &info)
	return true
}

func init() {
	chart.RegisterProcessFunc(Process)
}
