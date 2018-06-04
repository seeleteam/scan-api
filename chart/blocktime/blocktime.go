/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package blocktime

import (
	"scan-api/database"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	secondsInOneDay = float64(24 * 60 * 60)
)

//Process set an timer to calculate block time every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldBlockAvgTime()
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		ProcessOneDayBlockAvgTime(next)
	}
}

//ProcessOldBlockAvgTime Calculate average block time in the past days
func ProcessOldBlockAvgTime() {
	now := time.Now()
	//lastZeroTime := now.Add(-time.Hour * 24)
	todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
		_, err := database.GetOneDayBlockAvgTime(lastZeroTime.Unix())
		if err != mgo.ErrNotFound {
			break
		}
		if !ProcessOneDayBlockAvgTime(todayZeroTime) {
			break
		}
		todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
	}
}

//ProcessOneDayBlockAvgTime Calculate the average time for all blocks of a day
func ProcessOneDayBlockAvgTime(day time.Time) bool {
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

	var info database.DBOneDayBlockAvgTime

	if len(dbBlocks) > 0 {
		info.AvgTime = secondsInOneDay / float64(len(dbBlocks))
	} else {
		info.AvgTime = 0.0
	}

	info.TimeStamp = lastDayZeroTime.Unix()
	database.AddOneDayBlockAvgTime(&info)
	return true
}
