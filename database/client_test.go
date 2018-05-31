/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */
package database

import (
	"scan-api/log"
	"sort"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
)

func TestOnGetBlockByHeight(t *testing.T) {
	err := removeBlock(1000000)
	if err != nil && err != mgo.ErrNotFound {
		t.Errorf("db error, %v", err)
	}

	testBlock := DBBlock{
		HeadHash:        "0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57a5e9aa",
		PreHash:         "0x5ea58e84f6740d91972e5b77b96bc35f8ba6fe229782fe10497e6e88bac5c0e2",
		Height:          1000000,
		StateHash:       "",
		Timestamp:       1526867297,
		Difficulty:      "10000000",
		TotalDifficulty: "20000000",
		Creator:         "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Nonce:           "11793544988386285568",
		TxHash:          "",
	}

	err = AddBlock(&testBlock)
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	dbBlock, err := GetBlockByHeight(1000000)
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	if testBlock.HeadHash != dbBlock.HeadHash ||
		testBlock.PreHash != dbBlock.PreHash ||
		testBlock.Height != dbBlock.Height ||
		testBlock.StateHash != dbBlock.StateHash ||
		testBlock.Timestamp != dbBlock.Timestamp ||
		testBlock.Difficulty != dbBlock.Difficulty ||
		testBlock.TotalDifficulty != dbBlock.TotalDifficulty ||
		testBlock.Creator != dbBlock.Creator ||
		testBlock.Nonce != dbBlock.Nonce ||
		testBlock.TxHash != dbBlock.TxHash {
		t.Errorf("data is not match")
	}
}

func TestOnGetBlockByHash(t *testing.T) {
	err := removeBlock(1000000)
	if err != nil && err != mgo.ErrNotFound {
		t.Errorf("db error, %v", err)
	}

	testBlock := DBBlock{
		HeadHash:        "0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57bbbbbb",
		PreHash:         "0x5ea58e84f6740d91972e5b77b96bc35f8ba6fe229782fe10497e6e88bac5c0e2",
		Height:          1000000,
		StateHash:       "",
		Timestamp:       1526867297,
		Difficulty:      "10000000",
		TotalDifficulty: "20000000",
		Creator:         "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Nonce:           "11793544988386285568",
		TxHash:          "",
	}

	err = AddBlock(&testBlock)
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	dbBlock, err := GetBlockByHash("0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57bbbbbb")
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	if testBlock.HeadHash != dbBlock.HeadHash ||
		testBlock.PreHash != dbBlock.PreHash ||
		testBlock.Height != dbBlock.Height ||
		testBlock.StateHash != dbBlock.StateHash ||
		testBlock.Timestamp != dbBlock.Timestamp ||
		testBlock.Difficulty != dbBlock.Difficulty ||
		testBlock.TotalDifficulty != dbBlock.TotalDifficulty ||
		testBlock.Creator != dbBlock.Creator ||
		testBlock.Nonce != dbBlock.Nonce ||
		testBlock.TxHash != dbBlock.TxHash {
		t.Errorf("data is not match")
	}
}

func TestOnGetBlocksByHeight(t *testing.T) {
	begin := 10000000
	end := 10000010

	for i := begin; i <= end; i++ {
		err := removeBlock(uint64(i))
		if err != nil && err != mgo.ErrNotFound {
			t.Errorf("db error, %v", err)
		}
	}

	for i := begin; i <= end; i++ {
		testBlock := DBBlock{
			HeadHash:        "0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57bbbbbb",
			PreHash:         "0x5ea58e84f6740d91972e5b77b96bc35f8ba6fe229782fe10497e6e88bac5c0e2",
			Height:          int64(i),
			StateHash:       "",
			Timestamp:       1526867297,
			Difficulty:      "10000000",
			TotalDifficulty: "20000000",
			Creator:         "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
			Nonce:           "11793544988386285568",
			TxHash:          "",
		}

		err := AddBlock(&testBlock)
		if err != nil {
			t.Errorf("db error, %v", err)
		}
	}

	dbBlocks, err := GetBlocksByHeight(uint64(begin), uint64(end))
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	if len(dbBlocks) != end-begin {
		t.Errorf("db error, get number error")
	}
}

func TestOnGetBlocksByTime(t *testing.T) {
	now := time.Now()
	now = now.Add(24 * 365 * time.Hour)
	begin := 10000000
	end := 10000010

	for i := begin; i <= end; i++ {
		err := removeBlock(uint64(i))
		if err != nil && err != mgo.ErrNotFound {
			t.Errorf("db error, %v", err)
		}
	}

	for i := begin; i <= end; i++ {
		blockTime := now.Add(-24 * time.Duration(i-begin) * time.Hour)
		testBlock := DBBlock{
			HeadHash:        "0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57bbbbbb",
			PreHash:         "0x5ea58e84f6740d91972e5b77b96bc35f8ba6fe229782fe10497e6e88bac5c0e2",
			Height:          int64(i),
			StateHash:       "",
			Timestamp:       blockTime.Unix(),
			Difficulty:      "10000000",
			TotalDifficulty: "20000000",
			Creator:         "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
			Nonce:           "11793544988386285568",
			TxHash:          "",
		}

		err := AddBlock(&testBlock)
		if err != nil {
			t.Errorf("db error, %v", err)
		}
	}

	endTime := now.Add(-24 * time.Duration(end-begin) * time.Hour)
	dbBlocks, err := GetBlocksByTime(endTime.Unix(), now.Unix())
	if err != nil {
		t.Errorf("db error, %v", err)
	}

	if len(dbBlocks) != end-begin+1 {
		t.Errorf("db error, get number error")
	}
}

func TestOnGetBlockHeight(t *testing.T) {
	dropCollection(blockTbl)

	begin := 1
	end := 100

	for i := begin; i <= end; i++ {
		err := removeBlock(uint64(i))
		if err != nil && err != mgo.ErrNotFound {
			t.Errorf("db error, %v", err)
		}
	}

	for i := begin; i <= end; i++ {
		testBlock := DBBlock{
			HeadHash:        "0x000000c67385d722011c158fc22e88733c083da69fdd721bb13a05bc57bbbbbb",
			PreHash:         "0x5ea58e84f6740d91972e5b77b96bc35f8ba6fe229782fe10497e6e88bac5c0e2",
			Height:          int64(i),
			StateHash:       "",
			Timestamp:       1526867297,
			Difficulty:      "10000000",
			TotalDifficulty: "20000000",
			Creator:         "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
			Nonce:           "11793544988386285568",
			TxHash:          "",
		}

		err := AddBlock(&testBlock)
		if err != nil {
			t.Errorf("db error, %v", err)
		}
	}

	height, err := GetBlockHeight()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if height != uint64(end) {
		t.Fatalf("db error, get number error")
	}
}

func TestOnAddTx(t *testing.T) {
	dropCollection(txTbl)

	err := removeTx(10000000)
	if err != nil && err != mgo.ErrNotFound {
		t.Fatalf("db error, %v", err)
	}

	testTx := DBTx{
		Hash:         "0x2919c60a1c1d98cac0d33d336761571f2724fc15e2d6ced5b002e35c80bbbbbb",
		From:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		To:           "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Amount:       "20000000000",
		AccountNonce: "0",
		Timestamp:    "1526867474273961984",
		Payload:      "",
		Block:        "99999",
		Idx:          10000000,
	}

	err = AddTx(&testTx)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbTx, err := GetTxByIdx(10000000)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if testTx != *dbTx {
		t.Fatalf("data is not match")
	}
}

func TestOnGetTxsByIdx(t *testing.T) {
	dropCollection(txTbl)

	begin := 0
	end := 100

	for i := begin; i <= end; i++ {
		err := removeTx(uint64(i))
		if err != nil && err != mgo.ErrNotFound {
			t.Fatalf("db error, %v", err)
		}
	}

	for i := begin; i <= end; i++ {
		testTx := DBTx{
			Hash:         "0x2919c60a1c1d98cac0d33d336761571f2724fc15e2d6ced5b002e35c80bbbbbb",
			From:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			To:           "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
			Amount:       "20000000000",
			AccountNonce: "0",
			Timestamp:    "1526867474273961984",
			Payload:      "",
			Block:        "99999",
			Idx:          int64(i),
		}

		err := AddTx(&testTx)
		if err != nil {
			t.Fatalf("db error, %v", err)
		}
	}

	dbTxs, err := GetTxsByIdx(uint64(begin), uint64(end))
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbTxs) != end-begin {
		t.Fatalf("db error, get number error")
	}
}

func TestOnGetTxByHash(t *testing.T) {
	dropCollection(txTbl)

	err := removeTx(10000000)
	if err != nil && err != mgo.ErrNotFound {
		t.Fatalf("db error, %v", err)
	}

	testTx := DBTx{
		Hash:         "0x2919c60a1c1d98cac0d33d336761571f2724fc15e2d6ced5b002e35c80bbbbbb",
		From:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		To:           "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Amount:       "20000000000",
		AccountNonce: "0",
		Timestamp:    "1526867474273961984",
		Payload:      "",
		Block:        "99999",
		Idx:          10000000,
	}

	err = AddTx(&testTx)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbTx, err := GetTxByHash("0x2919c60a1c1d98cac0d33d336761571f2724fc15e2d6ced5b002e35c80bbbbbb")
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if testTx != *dbTx {
		t.Fatalf("data is not match")
	}
}

func TestOnGetTxCnt(t *testing.T) {
	dropCollection(txTbl)

	begin := 0
	end := 100

	for i := begin; i < end; i++ {
		err := removeTx(uint64(i))
		if err != nil && err != mgo.ErrNotFound {
			t.Fatalf("db error, %v", err)
		}
	}

	for i := begin; i < end; i++ {
		testTx := DBTx{
			Hash:         "0x2919c60a1c1d98cac0d33d336761571f2724fc15e2d6ced5b002e35c80bbbbbb",
			From:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			To:           "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
			Amount:       "20000000000",
			AccountNonce: "0",
			Timestamp:    "1526867474273961984",
			Payload:      "",
			Block:        "99999",
			Idx:          int64(i),
		}

		err := AddTx(&testTx)
		if err != nil {
			t.Fatalf("db error, %v", err)
		}
	}

	count, err := GetTxCnt()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if count != uint64(end) {
		t.Fatalf("db error, get number error")
	}
}

func TestOnAddAccount(t *testing.T) {
	dropCollection(accTbl)

	testAccount := DBAccount{
		Address: "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Balance: 1000,
		TxCount: 1000,
		Mined:   3000,
	}

	err := AddAccount(&testAccount)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbAccount, err := GetAccountByAddress("0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd")
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if testAccount.Address != dbAccount.Address ||
		testAccount.Balance != dbAccount.Balance ||
		testAccount.TxCount != dbAccount.TxCount ||
		testAccount.Mined != dbAccount.Mined {
		t.Fatalf("data is not match")
	}
}

func TestOnUpdateAccount(t *testing.T) {
	dropCollection(accTbl)

	testAccount := DBAccount{
		Address: "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Balance: 1000,
		TxCount: 1000,
		Mined:   3000,
	}

	err := AddAccount(&testAccount)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbAccountTx := DBAccountTx{
		Hash:      "11111",
		Block:     100,
		From:      "aaaa",
		To:        "bbbb",
		Amount:    100,
		Timestamp: 1000,
		TxFee:     0,
		InOrOut:   false,
	}
	var txs []DBAccountTx
	txs = append(txs, dbAccountTx)
	err = UpdateAccount("0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 100, &txs)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbAccount, err := GetAccountByAddress("0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd")
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if dbAccount.Balance != 1000+100 {
		t.Fatal("data is not match")
	}

	if len(dbAccount.Txs) != 1 {
		t.Fatal("data is not match")
	}

	if dbAccount.Txs[0] != dbAccountTx {
		t.Fatal("data is not match")
	}
}

func TestOnUpdateAccountMinedBlock(t *testing.T) {
	dropCollection(accTbl)

	testAccount := DBAccount{
		Address: "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Balance: 1000,
		TxCount: 1000,
		Mined:   3000,
	}

	err := AddAccount(&testAccount)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	err = UpdateAccountMinedBlock("0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 100)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbAccount, err := GetAccountByAddress("0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd")
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if dbAccount.Mined != 3000+100 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetAccounts(t *testing.T) {
	dropCollection(accTbl)

	testAccount1 := DBAccount{
		Address: "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Balance: 1000,
		TxCount: 1000,
		Mined:   3000,
	}

	err := AddAccount(&testAccount1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	testAccount2 := DBAccount{
		Address: "0x4dd6881d13ab5152127fdsfsc5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd",
		Balance: 1020,
		TxCount: 1000,
		Mined:   3000,
	}

	err = AddAccount(&testAccount2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbAccounts, err := GetAccounts(2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbAccounts) != 2 {
		t.Fatalf("data is not match")
	}

	if dbAccounts[0].Balance != 1020 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayTransInfo(t *testing.T) {
	dropCollection(chartTxTbl)

	oneDayTxInfo := DBOneDayTxInfo{
		TotalTxs:    2000,
		TotalBlocks: 1000,
		TimeStamp:   1527350400,
	}

	err := AddOneDayTransInfo(&oneDayTxInfo)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayTxInfo, err := GetOneDayTransInfo(1527350400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayTxInfo != *dbOneDayTxInfo {
		t.Fatalf("data is not match")
	}
}

func TestOnGetTransInfoChart(t *testing.T) {
	dropCollection(chartTxTbl)

	oneDayTxInfo1 := DBOneDayTxInfo{
		TotalTxs:    2000,
		TotalBlocks: 1000,
		TimeStamp:   1527350400,
	}

	err := AddOneDayTransInfo(&oneDayTxInfo1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayTxInfo2 := DBOneDayTxInfo{
		TotalTxs:    2000,
		TotalBlocks: 1000,
		TimeStamp:   1527350400,
	}

	err = AddOneDayTransInfo(&oneDayTxInfo2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	DBOneDayTxInfos, err := GetTransInfoChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(DBOneDayTxInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayHashRate(t *testing.T) {
	dropCollection(chartHashRateTbl)

	oneDayHashRate := DBOneDayHashRate{
		HashRate:  1000,
		TimeStamp: 1527350400,
	}

	err := AddOneDayHashRate(&oneDayHashRate)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayHashRate, err := GetOneDayHashRate(1527350400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayHashRate != *dbOneDayHashRate {
		t.Fatalf("data is not match")
	}
}

func TestOnGetHashRateChart(t *testing.T) {
	dropCollection(chartHashRateTbl)

	oneDayHashRate1 := DBOneDayHashRate{
		HashRate:  1000,
		TimeStamp: 1527350400,
	}

	err := AddOneDayHashRate(&oneDayHashRate1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayHashRate2 := DBOneDayHashRate{
		HashRate:  2000,
		TimeStamp: 1527350410,
	}

	err = AddOneDayHashRate(&oneDayHashRate2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayHashRateInfos, err := GetHashRateChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbOneDayHashRateInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlockDifficulty(t *testing.T) {
	dropCollection(chartBlockDifficultyTbl)

	oneDayBlockDifficulty := DBOneDayBlockDifficulty{
		Difficulty: 1200.98,
		TimeStamp:  1527350400,
	}

	err := AddOneDayBlockDifficulty(&oneDayBlockDifficulty)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockDifficulty, err := GetOneDayBlockDifficulty(1527350400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayBlockDifficulty != *dbOneDayBlockDifficulty {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlockDifficultyChart(t *testing.T) {
	dropCollection(chartBlockDifficultyTbl)

	oneDayBlockDifficulty1 := DBOneDayBlockDifficulty{
		Difficulty: 1200.98,
		TimeStamp:  1527350400,
	}

	err := AddOneDayBlockDifficulty(&oneDayBlockDifficulty1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayBlockDifficulty2 := DBOneDayBlockDifficulty{
		Difficulty: 1200.98,
		TimeStamp:  1528956400,
	}

	err = AddOneDayBlockDifficulty(&oneDayBlockDifficulty2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockDifficultyInfos, err := GetOneDayBlockDifficultyChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbOneDayBlockDifficultyInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlockAvgTime(t *testing.T) {
	dropCollection(chartBlockAvgTimeTbl)

	oneDayBlockAvgTime := DBOneDayBlockAvgTime{
		AvgTime:   100.6,
		TimeStamp: 1527350400,
	}

	err := AddOneDayBlockAvgTime(&oneDayBlockAvgTime)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockAvgTime, err := GetOneDayBlockAvgTime(1527350400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayBlockAvgTime != *dbOneDayBlockAvgTime {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlockAvgTimeChart(t *testing.T) {
	dropCollection(chartBlockAvgTimeTbl)

	oneDayBlockAvgTime1 := DBOneDayBlockAvgTime{
		AvgTime:   100.6,
		TimeStamp: 1527350400,
	}

	err := AddOneDayBlockAvgTime(&oneDayBlockAvgTime1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayBlockAvgTime2 := DBOneDayBlockAvgTime{
		AvgTime:   100.6,
		TimeStamp: 1527879400,
	}

	err = AddOneDayBlockAvgTime(&oneDayBlockAvgTime2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockAvgTimeInfos, err := GetOneDayBlockAvgTimeChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbOneDayBlockAvgTimeInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlock(t *testing.T) {
	dropCollection(chartBlockTbl)

	oneDayBlockInfo := DBOneDayBlockInfo{
		TotalBlocks: 100,
		Rewards:     200,
		TimeStamp:   1527879400,
	}

	err := AddOneDayBlock(&oneDayBlockInfo)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockInfo, err := GetOneDayBlock(1527879400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayBlockInfo != *dbOneDayBlockInfo {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayBlocksChart(t *testing.T) {
	dropCollection(chartBlockTbl)

	oneDayBlockInfo1 := DBOneDayBlockInfo{
		TotalBlocks: 100,
		Rewards:     200,
		TimeStamp:   1527879400,
	}

	err := AddOneDayBlock(&oneDayBlockInfo1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayBlockInfo2 := DBOneDayBlockInfo{
		TotalBlocks: 100,
		Rewards:     200,
		TimeStamp:   1527969400,
	}

	err = AddOneDayBlock(&oneDayBlockInfo2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayBlockInfos, err := GetOneDayBlocksChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbOneDayBlockInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayAddress(t *testing.T) {
	dropCollection(chartAddressTbl)

	oneDayAddressInfo := DBOneDayAddressInfo{
		TotalAddresss: 100,
		TodayIncrease: 50,
		TimeStamp:     1527969400,
	}

	err := AddOneDayAddress(&oneDayAddressInfo)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayAddressInfo, err := GetOneDayAddress(1527969400)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if oneDayAddressInfo != *dbOneDayAddressInfo {
		t.Fatalf("data is not match")
	}
}

func TestOnGetOneDayAddressesChart(t *testing.T) {
	dropCollection(chartAddressTbl)

	oneDayAddressInfo1 := DBOneDayAddressInfo{
		TotalAddresss: 100,
		TodayIncrease: 50,
		TimeStamp:     1527969400,
	}

	err := AddOneDayAddress(&oneDayAddressInfo1)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	oneDayAddressInfo2 := DBOneDayAddressInfo{
		TotalAddresss: 200,
		TodayIncrease: 30,
		TimeStamp:     1526769400,
	}

	err = AddOneDayAddress(&oneDayAddressInfo2)
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	dbOneDayAddressInfos, err := GetOneDayAddressesChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(dbOneDayAddressInfos) != 2 {
		t.Fatalf("data is not match")
	}
}

//MinerRankInfoSlice rank array
type MinerRankInfoSlice []DBSingleMinerRankInfo

func (s MinerRankInfoSlice) Len() int           { return len(s) }
func (s MinerRankInfoSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MinerRankInfoSlice) Less(i, j int) bool { return s[i].Mined > s[j].Mined }
func TestOnGetTopMinerChart(t *testing.T) {
	RemoveTopMinerInfo()

	var TopSevenDaysMiners MinerRankInfoSlice

	miner := DBSingleMinerRankInfo{
		Address:    "11111",
		Mined:      100,
		Percentage: 100 / (100 + 200 + 300),
	}
	TopSevenDaysMiners = append(TopSevenDaysMiners, miner)

	miner2 := DBSingleMinerRankInfo{
		Address:    "22222",
		Mined:      100,
		Percentage: 100 / (100 + 200 + 300),
	}
	TopSevenDaysMiners = append(TopSevenDaysMiners, miner2)

	miner3 := DBSingleMinerRankInfo{
		Address:    "33333",
		Mined:      100,
		Percentage: 100 / (100 + 200 + 300),
	}
	TopSevenDaysMiners = append(TopSevenDaysMiners, miner3)

	sort.Stable(TopSevenDaysMiners)
	dbRank := DBMinerRankInfo{Rank: TopSevenDaysMiners}

	AddTopMinerInfo(&dbRank)

	topMiners, err := GetTopMinerChart()
	if err != nil {
		t.Fatalf("db error, %v", err)
	}

	if len(topMiners[0].Rank) != 3 {
		t.Fatalf("data is not match")
	}
}

func init() {
	log.NewLogger("debug", false)
}
