/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/seeleteam/scan-api/common"
	"github.com/seeleteam/scan-api/database"
)

var db = &common.DataBaseConfig{
	DataBaseName:     "seele",
	DataBaseMode:     "single",
	DataBaseConnURLs: []string{"127.0.0.1:27017"},
}

func Benchmark_LastBlock(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		lastBlockHeight, lastBlockTime, err := dbClient.GetBlockProTime()
		if err == nil {
			createRetLastblockInfo(lastBlockHeight, lastBlockTime)
		}
	}
}

func Benchmark_txcount(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := dbClient.GetTxCnt()
		if err != nil {
			b.Errorf("error:%s", err)
		}
	}
}

func Benchmark_blockTxsTps(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := dbClient.GetBlockTxsTps()
		if err != nil {
			b.Errorf("error:%s", err)
		}
	}
}

func Benchmark_Txstat(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		now := time.Now()
		startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
		todayDate := now.Format("2006-01-02")
		txs, err := dbClient.GetTxHis(startDate, todayDate)
		if err != nil {
			b.Errorf("error:%s", err)
			return
		}
		for _, tx := range txs {
			ymd := strings.Split(tx.Stime, "-")
			year, _ := strconv.Atoi(ymd[0])
			month, _ := strconv.Atoi(ymd[1])
			day, _ := strconv.Atoi(ymd[2])
			dateTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
			timestamp := dateTime.Unix()
			tx.Stime = strconv.FormatInt(timestamp, 10)
		}
	}
}

func Benchmark_accountcount(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := dbClient.GetAccountCnt()
		if err != nil {
			b.Errorf("error:%s", err)
			return
		}
	}
}

func Benchmark_contractcount(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := dbClient.GetContractCnt()
		if err != nil {
			b.Errorf("error:%s", err)
			return
		}
	}
}
