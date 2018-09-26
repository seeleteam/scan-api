/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/rpc"
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
	ShardNumber int    `json:"shardnumber"`
	Height      uint64 `json:"height"`
	Age         string `json:"age"`
	Txn         int    `json:"txn"`
	Miner       string `json:"miner"`
	Reward      int64  `json:"reward"`
	Fee         int64  `json:"fee"`
}

//RetDetailBlockInfo describle the block info in the block detail page which send to the frontend
type RetDetailBlockInfo struct {
	ShardNumber int      `json:"shardnumber"`
	HeadHash    string   `json:"headHash"`
	PreHash     string   `json:"preBlockHash"`
	Height      uint64   `json:"height"`
	Age         string   `json:"age"`
	Difficulty  *big.Int `json:"difficulty"`
	Miner       string   `json:"miner"`
	Nonce       string   `json:"nonce"`
	TxCount     int      `json:"txcount"`

	MaxHeight uint64 `json:"maxheight"`
	MinHeight uint64 `json:"minheight"`
}

//CountsTime
type CountsTime struct {
	Stime   string `json:"stime"`
	TxCount int64  `json:"txcount"`
}

//Lastblock
type Lastblock struct {
	LastblockHeight int64 `json:"lastblockHeight"`
	LastblockTime   int64 `json:"lastblockTime"`
}

//RetSimpleTxInfo describle the transaction info in the transaction detail page which send to the frontend
type RetSimpleTxInfo struct {
	TxType      int         `json:"txtype"`
	ShardNumber int         `json:"shardnumber"`
	TxHash      string      `json:"txHash"`
	Block       uint64      `json:"block"`
	Age         string      `json:"age"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Value       int64       `json:"value"`
	Pending     bool        `json:"pending"`
	Fee         int64       `json:"fee"`
	Receipt     rpc.Receipt `json:"receipt"`
	Nonce       string      `json:"nonce"`
}

//RetDetailTxInfo describle the transaction detail info in the transaction detail page which send to the frontend
type RetDetailTxInfo struct {
	TxType       int         `json:"txtype"`
	ShardNumber  int         `json:"shardnumber"`
	TxHash       string      `json:"txHash"`
	Block        uint64      `json:"block"`
	Age          string      `json:"age"`
	From         string      `json:"from"`
	To           string      `json:"to"`
	Value        int64       `json:"value"`
	Pending      bool        `json:"pending"`
	Fee          int64       `json:"fee"`
	AccountNonce string      `json:"accountNonce"`
	Payload      string      `json:"payload"`
	Receipt      rpc.Receipt `json:"receipt"`
}

//RetSimpleAccountInfo describle the account info in the account list page which send to the frontend
type RetSimpleAccountInfo struct {
	AccType     int     `json:"accType"`
	ShardNumber int     `json:"shardnumber"`
	Rank        int     `json:"rank"`
	Address     string  `json:"address"`
	Balance     int64   `json:"balance"`
	Percentage  float64 `json:"percentage"`
	TxCount     int64   `json:"txcount"`
}

//RetSimpleAccountHome
type RetSimpleAccountHome struct {
	Number     int    `json:"number"`
	Address    string `json:"address"`
	Balance    int64  `json:"balance"`
	Percentage string `json:"percentage"`
}

//RetDetailAccountTxInfo describle the tx info contained by the RetDetailAccountInfo
type RetDetailAccountTxInfo struct {
	ShardNumber int    `json:"shardnumber"`
	TxType      int    `json:"txtype"`
	Hash        string `json:"hash"`
	Block       uint64 `json:"block"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       int64  `json:"value"`
	Age         string `json:"age"`
	Fee         int64  `json:"fee"`
	InOrOut     bool   `json:"inorout"`
	Pending     bool   `json:"pending"`
}

//RetDetailAccountInfo describle the detail account info which send to the frontend
type RetDetailAccountInfo struct {
	AccType              int                      `json:"accType"`
	ShardNumber          int                      `json:"shardnumber"`
	Address              string                   `json:"address"`
	Balance              int64                    `json:"balance"`
	Percentage           float64                  `json:"percentage"`
	TxCount              int64                    `json:"txcount"`
	ContractCreationCode string                   `json:"contractCreationCode"`
	Txs                  []RetDetailAccountTxInfo `json:"txs"`
}

//createRetinfor
func createRetinfor(data int64, theTime string) *CountsTime {
	var ret CountsTime
	ret.TxCount = data
	ret.Stime = theTime
	return &ret
}

//createRetLastblockInfo converts the given dbblock to the Lastblock
func createRetLastblockInfo(lastblockHeight int64, lastblockTime int64) *Lastblock {
	var ret Lastblock
	ret.LastblockHeight = lastblockHeight
	ret.LastblockTime = lastblockTime
	return &ret
}

//createRetSimpleBlockInfo converts the given dbblock to the retsimpleblockinfo
func createRetSimpleBlockInfo(blockInfo *database.DBBlock) *RetSimpleBlockInfo {
	var ret RetSimpleBlockInfo
	var blockFee int64
	ret.Miner = blockInfo.Creator
	ret.Height = uint64(blockInfo.Height)
	ret.Txn = len(blockInfo.Txs)
	timeStamp := big.NewInt(blockInfo.Timestamp)
	ret.Age = getElpasedTimeDesc(timeStamp)
	ret.ShardNumber = blockInfo.ShardNumber
	ret.Reward = blockInfo.Reward

	for i := 0; i < len(blockInfo.Txs); i++ {
		blockFee += blockInfo.Txs[i].Fee
	}
	ret.Fee = blockFee
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
	ret.ShardNumber = blockInfo.ShardNumber
	return &ret
}

//createRetSimpleTxInfo converts the given dbtx to the retsimpletxinfo
func createRetSimpleTxInfo(transaction *database.DBTx) *RetSimpleTxInfo {
	var ret RetSimpleTxInfo
	ret.TxType = transaction.TxType
	ret.TxHash = transaction.Hash
	ret.Block = transaction.Block
	ret.From = transaction.From
	ret.To = transaction.To
	ret.Value = transaction.Amount
	ret.Pending = transaction.Pending
	ret.Fee = transaction.Fee
	ret.Nonce = transaction.AccountNonce
	timeStamp := big.NewInt(0)
	if timeStamp.UnmarshalText([]byte(transaction.Timestamp)) == nil {
		ret.Age = getElpasedTimeDesc(timeStamp)
	}
	ret.ShardNumber = transaction.ShardNumber
	ret.Receipt = transaction.Receipt
	return &ret
}

func createRetDetailTxInfo(transaction *database.DBTx) *RetDetailTxInfo {
	var ret RetDetailTxInfo
	ret.TxType = transaction.TxType
	ret.TxHash = transaction.Hash
	ret.Block = transaction.Block
	ret.From = transaction.From
	ret.To = transaction.To
	ret.Value = transaction.Amount
	ret.Pending = transaction.Pending
	ret.Fee = transaction.Fee
	timeStamp := big.NewInt(0)
	if timeStamp.UnmarshalText([]byte(transaction.Timestamp)) == nil {
		ret.Age = getElpasedTimeDesc(timeStamp)
	}
	ret.ShardNumber = transaction.ShardNumber
	ret.AccountNonce = transaction.AccountNonce
	ret.Payload = transaction.Payload
	ret.Receipt = transaction.Receipt
	return &ret
}

//createRetSimpleAccountInfo converts the given dbaccount to the retsimpleaccountinfo
func createRetSimpleAccountInfo(account *database.DBAccount, ttBalance int64) *RetSimpleAccountInfo {
	var ret RetSimpleAccountInfo
	ret.AccType = account.AccType
	ret.Address = account.Address
	ret.Balance = account.Balance
	ret.TxCount = account.TxCount
	ret.Percentage = (float64(ret.Balance) / float64(ttBalance))
	ret.ShardNumber = account.ShardNumber
	return &ret
}

//createHomeRetSimpleAccountInfo converts the given dbaccount to the RetSimpleAccountHome
func createHomeRetSimpleAccountInfo(account *database.DBAccount) *RetSimpleAccountHome {
	var ret RetSimpleAccountHome
	ret.Address = account.Address
	ret.Balance = account.Balance
	AccountBalance := strconv.FormatFloat((float64(ret.Balance) / 1000000000), 'f', -1, 64)
	ret.Percentage = AccountBalance
	return &ret
}

//createRetDetailAccountInfo converts the given dbaccount to the tetdetailaccountInfo
func createRetDetailAccountInfo(account *database.DBAccount, txs []*database.DBTx, ttBalance int64) *RetDetailAccountInfo {
	var ret RetDetailAccountInfo
	ret.AccType = account.AccType
	ret.Address = account.Address
	ret.Balance = account.Balance
	ret.TxCount = account.TxCount
	ret.Percentage = (float64(ret.Balance) / float64(ttBalance))

	for i := 0; i < len(txs); i++ {
		var tx RetDetailAccountTxInfo
		tx.TxType = txs[i].TxType
		tx.Value = txs[i].Amount
		tx.Block = txs[i].Block
		tx.From = txs[i].From
		tx.Hash = txs[i].Hash
		tx.To = txs[i].To
		tx.ShardNumber = txs[i].ShardNumber
		if tx.From == account.Address {
			tx.InOrOut = false
		} else {
			tx.InOrOut = true
		}
		timeStamp := big.NewInt(0)
		if timeStamp.UnmarshalText([]byte(txs[i].Timestamp)) == nil {
			tx.Age = getElpasedTimeDesc(timeStamp)
		}

		tx.Fee = txs[i].Fee
		tx.Pending = txs[i].Pending
		ret.Txs = append(ret.Txs, tx)

		if txs[i].TxType == 1 {
			ret.ContractCreationCode = txs[i].Payload
		}
	}
	ret.ShardNumber = account.ShardNumber

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
