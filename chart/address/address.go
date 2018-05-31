/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package address

import (
	"scan-api/database"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	nullAddress = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
)

//Process set an timer to process address data every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessOldAddresses()
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 1, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		ProcessOneDayAddresses(next)
	}
}

//ProcessOldAddresses Process address data in the past days
func ProcessOldAddresses() {
	secondBlock, err := database.GetBlockByHeight(1)
	if err != nil {
		return
	}

	now := time.Now()
	lastZeroTime := now.Add(-time.Hour * 24)
	lastZeroTime = time.Date(lastZeroTime.Year(), lastZeroTime.Month(), lastZeroTime.Day(), 0, 0, 0, 0, lastZeroTime.Location())
	for {
		_, err := database.GetOneDayAddress(lastZeroTime.Unix())
		if err != mgo.ErrNotFound {
			break
		}
		if lastZeroTime.Unix() < secondBlock.Timestamp {
			lastZeroTime = lastZeroTime.Add(time.Hour * 24)
			break
		}
		lastZeroTime = lastZeroTime.Add(-time.Hour * 24)
	}

	beginTime := lastZeroTime
	for {
		if beginTime.Unix() > now.Unix() {
			break
		}

		if !ProcessOneDayAddresses(beginTime) {
			break
		}

		beginTime = beginTime.Add(time.Hour * 24)
	}

}

//ProcessOneDayAddresses Count all addresses within one day
func ProcessOneDayAddresses(day time.Time) bool {
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

	var info database.DBOneDayAddressInfo

	for i := 0; i < len(dbBlocks); i++ {
		for j := 0; j < len(dbBlocks[i].Txs); j++ {
			tx := dbBlocks[i].Txs[j]
			fromAddress := tx.From
			toAddress := tx.To
			if fromAddress != nullAddress {
				_, err := database.GetOneDaySingleAddressInfo(fromAddress)
				if err == mgo.ErrNotFound {
					info.TodayIncrease++
					fromAccount := &database.DBOneDaySingleAddressInfo{
						Address:   fromAddress,
						TimeStamp: lastDayZeroTime.Unix(),
					}
					database.AddOneDaySingleAddressInfo(fromAccount)
				}
			}

			_, err := database.GetOneDaySingleAddressInfo(toAddress)
			if err == mgo.ErrNotFound {
				info.TodayIncrease++
				toAccount := &database.DBOneDaySingleAddressInfo{
					Address:   toAddress,
					TimeStamp: lastDayZeroTime.Unix(),
				}
				database.AddOneDaySingleAddressInfo(toAccount)
			}
		}
	}

	last2DayZeroTime := lastDayZeroTime.Add(-time.Hour * 24)
	lastDayInfo, err := database.GetOneDayAddress(last2DayZeroTime.Unix())
	if err != mgo.ErrNotFound {
		info.TotalAddresss = lastDayInfo.TotalAddresss + info.TodayIncrease
	} else {
		info.TotalAddresss = info.TodayIncrease
	}
	info.TimeStamp = lastDayZeroTime.Unix()
	database.AddOneDayAddress(&info)
	return true
}
