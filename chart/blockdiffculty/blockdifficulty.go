/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package blockdifficulty

import (
	"scan-api/database"
	"strconv"
	"sync"
	"time"

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
		ProcessOneDayBlockDifficulty(next)
	}
}

//ProcessOldBlockDifficulty Calculate block difficulty in the past days
func ProcessOldBlockDifficulty() {
	now := time.Now()
	//lastZeroTime := now.Add(-time.Hour * 24)
	todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
		_, err := database.GetOneDayBlockDifficulty(lastZeroTime.Unix())
		if err != mgo.ErrNotFound {
			break
		}
		if !ProcessOneDayBlockDifficulty(todayZeroTime) {
			break
		}
		todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
	}
}

//ProcessOneDayBlockDifficulty Calculate the average difficulty of all blocks within one day
func ProcessOneDayBlockDifficulty(day time.Time) bool {
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
	database.AddOneDayBlockDifficulty(&info)
	return true
}
