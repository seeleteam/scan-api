/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package hashrate

import (
	"strconv"
	"sync"
	"time"

	"github.com/seeleteam/scan-api/chart"
	"github.com/seeleteam/scan-api/database"

	mgo "gopkg.in/mgo.v2"
)

const (
	secondsInOneDay = float64(24 * 60 * 60)
)

//Process set an timer to calculate hashrate every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldHashRate()
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		for i := 1; i <= chart.ShardCount; i++ {
			ProcessOneDayHashRate(i, next)
		}
	}
}

//ProcessOldHashRate Calculate hashrate in the past days
func ProcessOldHashRate() {
	for i := 1; i <= chart.ShardCount; i++ {
		now := time.Now()
		todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		for {

			lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
			_, err := chart.GChartDB.GetOneDayHashRate(i, lastZeroTime.Unix())
			if err != mgo.ErrNotFound {
				break
			}
			if !ProcessOneDayHashRate(i, todayZeroTime) {
				break
			}
			todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
		}
	}
}

//ProcessOneDayHashRate  Calculate the hashrate of a day
func ProcessOneDayHashRate(shardNumber int, day time.Time) bool {
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

	var info database.DBOneDayHashRate

	var difficulty uint64

	for i := 0; i < len(dbBlocks); i++ {
		diff, err := strconv.ParseUint(dbBlocks[i].Difficulty, 10, 64)
		if err == nil {
			difficulty += diff
		}
	}

	info.HashRate = (float64(difficulty) / float64(secondsInOneDay))
	info.TimeStamp = lastDayZeroTime.Unix()
	chart.GChartDB.AddOneDayHashRate(shardNumber, &info)
	return true
}

func init() {
	chart.RegisterProcessFunc(Process)
}
