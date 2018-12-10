/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"fmt"
	"testing"

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
