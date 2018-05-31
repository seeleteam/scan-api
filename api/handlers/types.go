/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"fmt"
	"math/big"
	"scan-api/database"
	"strconv"

	"time"
)

const (
	year  = 365 * 24 * 60 * 60
	month = 30 * 24 * 60 * 60
	day   = 24 * 60 * 60
	hour  = 60 * 60
	min   = 60
)

//RetSimpleBlockInfo describle the block info in the block list which send to the frontend
type RetSimpleBlockInfo struct {
	Height uint64 `json:"height"`
	Age    string `json:"age"`
	Txn    int    `json:"txn"`
	Miner  string `json:"miner"`
}

//RetDetailBlockInfo describle the block info in the block detail page which send to the frontend
type RetDetailBlockInfo struct {
	HeadHash   string   `json:"headHash"`
	PreHash    string   `json:"preBlockHash"`
	Height     uint64   `json:"height"`
	Age        string   `json:"age"`
	Difficulty *big.Int `json:"difficulty"`
	Miner      string   `json:"miner"`
	Nonce      string   `json:"nonce"`
	TxCount    int      `json:"txcount"`

	MaxHeight uint64 `json:"maxheight"`
	MinHeight uint64 `json:"minheight"`
}

//RetSimpleTxInfo describle the transaction info in the transaction detail page which send to the frontend
type RetSimpleTxInfo struct {
	TxHash string `json:"txHash"`
	Block  uint64 `json:"block"`
	Age    string `json:"age"`
	From   string `json:"from"`
	To     string `json:"to"`
	Value  string `json:"value"`
}

//RetSimpleAccountInfo describle the account info in the account list page which send to the frontend
type RetSimpleAccountInfo struct {
	Rank       int     `json:"rank"`
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
	Percentage float64 `json:"percentage"`
	TxCount    int     `json:"txcount"`
}

//RetDetailAccountTxInfo describle the tx info contained by the RetDetailAccountInfo
type RetDetailAccountTxInfo struct {
	Hash    string  `json:"hash"`
	Block   int64   `json:"block"`
	From    string  `json:"from"`
	To      string  `json:"to"`
	Amount  int64   `json:"amount"`
	Age     string  `json:"age"`
	TxFee   float64 `json:"txfee"`
	InOrOut bool    `json:"inorout"`
}

//RetDetailAccountInfo describle the detail account info which send to the frontend
type RetDetailAccountInfo struct {
	Address    string                   `json:"address"`
	Balance    float64                  `json:"balance"`
	Percentage float64                  `json:"percentage"`
	TxCount    int                      `json:"txcount"`
	Txs        []RetDetailAccountTxInfo `json:"txs"`
}

//createRetSimpleBlockInfo converts the given dbblock to the retsimpleblockinfo
func createRetSimpleBlockInfo(blockInfo *database.DBBlock) *RetSimpleBlockInfo {
	var ret RetSimpleBlockInfo
	ret.Miner = blockInfo.Creator
	ret.Height = uint64(blockInfo.Height)
	ret.Txn = len(blockInfo.Txs)
	timeStamp := big.NewInt(blockInfo.Timestamp)
	ret.Age = getElpasedTimeDesc(timeStamp)

	return &ret
}

//createRetDetailBlockInfo converts the given dbblock to the retdetailblockinfo
func createRetDetailBlockInfo(blockInfo *database.DBBlock, maxHeight, minHeight uint64) *RetDetailBlockInfo {
	var ret RetDetailBlockInfo
	ret.HeadHash = blockInfo.HeadHash
	ret.PreHash = blockInfo.PreHash
	ret.Height = uint64(blockInfo.Height)
	timeStamp := big.NewInt(blockInfo.Timestamp)
	ret.Age = getElpasedTimeDesc(timeStamp)

	difficulty := big.NewInt(0)
	if difficulty.UnmarshalText([]byte(blockInfo.Difficulty)) == nil {
		ret.Difficulty = difficulty
	}

	ret.Miner = blockInfo.Creator

	ret.Nonce = blockInfo.Nonce
	ret.TxCount = len(blockInfo.Txs)
	ret.MaxHeight = maxHeight
	ret.MinHeight = minHeight

	return &ret
}

//createRetSimpleTxInfo converts the given dbtx to the retsimpletxinfo
func createRetSimpleTxInfo(transaction *database.DBTx) *RetSimpleTxInfo {
	var ret RetSimpleTxInfo
	ret.TxHash = transaction.Hash
	ret.Block, _ = strconv.ParseUint(transaction.Block, 10, 64)
	ret.From = transaction.From
	ret.To = transaction.To
	ret.Value = transaction.Amount
	timeStamp := big.NewInt(0)
	if timeStamp.UnmarshalText([]byte(transaction.Timestamp)) == nil {
		ret.Age = getElpasedTimeDesc(timeStamp.Div(timeStamp, big.NewInt(1e9)))
	}
	return &ret
}

//createRetSimpleAccountInfo converts the given dbaccount to the retsimpleaccountinfo
func createRetSimpleAccountInfo(account *database.DBAccount) *RetSimpleAccountInfo {
	var ret RetSimpleAccountInfo
	ret.Address = account.Address
	ret.Balance = account.Balance
	ret.TxCount = account.TxCount
	ret.Percentage = 0.0
	return &ret
}

//createRetDetailAccountInfo converts the given dbaccount to the tetdetailaccountInfo
func createRetDetailAccountInfo(account *database.DBAccount) *RetDetailAccountInfo {
	var ret RetDetailAccountInfo
	ret.Address = account.Address
	ret.Balance = account.Balance
	ret.TxCount = account.TxCount
	ret.Percentage = 0.0
	for i := 0; i < len(account.Txs); i++ {
		var tx RetDetailAccountTxInfo
		tx.Amount = account.Txs[i].Amount
		tx.Block = account.Txs[i].Block
		tx.From = account.Txs[i].From
		tx.Hash = account.Txs[i].Hash
		tx.To = account.Txs[i].To
		tx.InOrOut = account.Txs[i].InOrOut
		tx.Age = getElpasedTimeDesc(big.NewInt(account.Txs[i].Timestamp))
		tx.TxFee = 0.0
		ret.Txs = append(ret.Txs, tx)
	}
	return &ret
}

//getElpasedTimeDesc Get the elapsed time from then until now
func getElpasedTimeDesc(t *big.Int) string {
	curTimeStamp := time.Now().Unix()
	minerTimeStamp := t.Int64()
	elpasedSec := curTimeStamp - minerTimeStamp
	switch {
	case elpasedSec > year:
		nYears := elpasedSec / year
		return fmt.Sprintf("%d years ago", nYears)
	case elpasedSec > month:
		nMonths := elpasedSec / month
		return fmt.Sprintf("%d months ago", nMonths)
	case elpasedSec > day:
		nDays := elpasedSec / day
		return fmt.Sprintf("%d days ago", nDays)
	case elpasedSec > hour:
		nHours := elpasedSec / hour
		return fmt.Sprintf("%d hours ago", nHours)
	case elpasedSec > min:
		nMins := elpasedSec / min
		return fmt.Sprintf("%d mins ago", nMins)
	default:
		nSecs := elpasedSec
		if nSecs <= 0 {
			nSecs = 1
		}
		return fmt.Sprintf("%d secs ago", nSecs)
	}

}

//RetOneDayTxInfo describle the transaction info in the transaction history chart page which send to the frontend
type RetOneDayTxInfo struct {
	TotalTxs      int
	TotalBlocks   int
	HashRate      float64
	Difficulty    float64
	AvgTime       float64
	Rewards       int64
	TotalAddresss int64
	TodayIncrease int64
	TimeStamp     int64
}
