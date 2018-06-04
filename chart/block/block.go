/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package block

import (
	"scan-api/database"
	"strconv"
	"sync"
	"time"

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
		ProcessOneDayBlocks(next)
	}
}

//ProcessOldBlocks Process block data mined in the past
func ProcessOldBlocks() {
	now := time.Now()
	//lastZeroTime := now.Add(-time.Hour * 24)
	todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		lastZeroTime := todayZeroTime.Add(-time.Hour * 24)
		_, err := database.GetOneDayBlock(lastZeroTime.Unix())
		if err != mgo.ErrNotFound {
			break
		}
		if !ProcessOneDayBlocks(todayZeroTime) {
			break
		}
		todayZeroTime = todayZeroTime.Add(-time.Hour * 24)
	}
}

//ProcessOneDayBlocks Process an single day block data
func ProcessOneDayBlocks(day time.Time) bool {
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

	var info database.DBOneDayBlockInfo

	for i := 0; i < len(dbBlocks); i++ {
		info.TotalBlocks++
		if len(dbBlocks[i].Txs) > 0 {
			amount, err := strconv.ParseInt(dbBlocks[i].Txs[0].Amount, 10, 64)
			if err == nil {
				info.Rewards += amount
			}
		}
	}

	info.TimeStamp = lastDayZeroTime.Unix()
	database.AddOneDayBlock(&info)
	return true
}
