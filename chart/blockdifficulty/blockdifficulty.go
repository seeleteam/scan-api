/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package blockdifficulty

import (
	"strconv"
	"sync"
	"time"

	"github.com/seeleteam/scan-api/chart"
	"github.com/seeleteam/scan-api/database"

	mgo "gopkg.in/mgo.v2"
)

//Process set an timer to process block difficulty every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldBlockDifficulty()

	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		for i := 1; i <= chart.ShardCount; i++ {
			ProcessOneDayBlockDifficulty(i, next)
		}
	}
}

//ProcessOldBlockDifficulty Calculate block difficulty in the past days
func ProcessOldBlockDifficulty() {
	for i := 1; i <= chart.ShardCount; i++ {
		now := time.Now()
		//lastZeroTime := now.Add(-time.Hour * 24)
		todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		for {

			lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
			_, err := chart.GChartDB.GetOneDayBlockDifficulty(i, lastZeroTime.Unix())
			if err != mgo.ErrNotFound {
				break
			}
			if !ProcessOneDayBlockDifficulty(i, todayZeroTime) {
				break
			}
			todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
		}
	}

}

//ProcessOneDayBlockDifficulty Calculate the average difficulty of all blocks within one day
func ProcessOneDayBlockDifficulty(shardNumber int, day time.Time) bool {
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

	var info database.DBOneDayBlockDifficulty
	info.Difficulty = 0.0

	for i := 0; i < len(dbBlocks); i++ {
		diff, err := strconv.ParseInt(dbBlocks[i].Difficulty, 10, 64)
		if err == nil {
			info.Difficulty += float64(diff)
		}
	}
	if len(dbBlocks) > 0 {
		info.Difficulty = info.Difficulty / float64(len(dbBlocks))
	}

	info.TimeStamp = lastDayZeroTime.Unix()
	chart.GChartDB.AddOneDayBlockDifficulty(shardNumber, &info)
	return true
}

func init() {
	chart.RegisterProcessFunc(Process)
}
