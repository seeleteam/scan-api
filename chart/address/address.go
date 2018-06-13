/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package address

import (
	"sync"
	"time"

	"github.com/seeleteam/scan-api/chart"
	"github.com/seeleteam/scan-api/database"
	mgo "gopkg.in/mgo.v2"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
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
		for i := 1; i <= chart.ShardCount; i++ {
			ProcessOneDayAddresses(i, next)
		}

	}
}

//ProcessOldAddresses Process address data in the past days
func ProcessOldAddresses() {
	for i := 1; i <= chart.ShardCount; i++ {
		secondBlock, err := chart.GChartDB.GetBlockByHeight(i, 1)
		if err != nil {
			continue
		}

		now := time.Now()
		lastZeroTime := now.Add(-time.Hour * 24)
		lastZeroTime = time.Date(lastZeroTime.Year(), lastZeroTime.Month(), lastZeroTime.Day(), 0, 0, 0, 0, lastZeroTime.Location())
		for {
			_, err := chart.GChartDB.GetOneDayAddress(i, lastZeroTime.Unix())
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

			if !ProcessOneDayAddresses(i, beginTime) {
				break
			}

			beginTime = beginTime.Add(time.Hour * 24)
		}

	}

}

//ProcessOneDayAddresses Count all addresses within one day
func ProcessOneDayAddresses(shardNumber int, day time.Time) bool {
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

	var info database.DBOneDayAddressInfo

	for i := 0; i < len(dbBlocks); i++ {
		for j := 0; j < len(dbBlocks[i].Txs); j++ {
			tx := dbBlocks[i].Txs[j]
			fromAddress := tx.From
			toAddress := tx.To
			if fromAddress != nullAddress {
				_, err := chart.GChartDB.GetOneDaySingleAddressInfo(shardNumber, fromAddress)
				if err == mgo.ErrNotFound {
					info.TodayIncrease++
					fromAccount := &database.DBOneDaySingleAddressInfo{
						Address:   fromAddress,
						TimeStamp: lastDayZeroTime.Unix(),
					}
					chart.GChartDB.AddOneDaySingleAddressInfo(shardNumber, fromAccount)
				}
			}

			_, err := chart.GChartDB.GetOneDaySingleAddressInfo(shardNumber, toAddress)
			if err == mgo.ErrNotFound {
				info.TodayIncrease++
				toAccount := &database.DBOneDaySingleAddressInfo{
					Address:   toAddress,
					TimeStamp: lastDayZeroTime.Unix(),
				}
				chart.GChartDB.AddOneDaySingleAddressInfo(shardNumber, toAccount)
			}
		}
	}

	last2DayZeroTime := lastDayZeroTime.Add(-time.Hour * 24)
	lastDayInfo, err := chart.GChartDB.GetOneDayAddress(shardNumber, last2DayZeroTime.Unix())
	if err != mgo.ErrNotFound {
		info.TotalAddresss = lastDayInfo.TotalAddresss + info.TodayIncrease
	} else {
		info.TotalAddresss = info.TodayIncrease
	}
	info.TimeStamp = lastDayZeroTime.Unix()
	chart.GChartDB.AddOneDayAddress(shardNumber, &info)
	return true
}

func init() {
	chart.RegisterProcessFunc(Process)
}
