/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"fmt"
	"testing"

	"github.com/seeleteam/scan-api/database"
)

func Benchmark_GetHomeAccounts(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		var Account []*RetSimpleAccountHome
		Accounts := dbClient.GetAccountsByHome()
		for i := 0; i < len(Accounts); i++ {
			data := Accounts[i]
			simpleTx := createHomeRetSimpleAccountInfo(data)
			Account = append(Account, simpleTx)
		}
	}
}

func Benchmark_GetMinerAccounts(b *testing.B) {
	dbClient := database.NewDBClient(db, 1)
	if dbClient == nil {
		fmt.Printf("init database error")
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := dbClient.GetMinerAccounts(MINERRANKSIZE)
		if err != nil {
			b.Errorf("error:%s", err)
		}
	}
}
