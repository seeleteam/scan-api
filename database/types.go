/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"strconv"
	"time"

	"github.com/seeleteam/scan-api/rpc"
)

//DBSimpleTxInBlock describle the transaction info contained by dbblock which stored in the database
type DBSimpleTxInBlock struct {
	Hash      string `bson:"hash"`
	From      string `bson:"from"`
	To        string `bson:"to"`
	Amount    int64  `bson:"amount"`
	Timestamp string `bson:"timestamp"`
	Fee       int64  `bson:"fee"`
}

type DBSimpleTxs struct {
	Stime   string `json:"stime"`
	TxCount int    `json:"txcount"`
}

//DBBlock describle the block info which stored in the database
type DBBlock struct {
	HeadHash        string                  `bson:"headHash"`
	PreHash         string                  `bson:"preBlockHash"`
	Height          int64                   `bson:"height"`
	StateHash       string                  `bson:"stateHash"`
	Timestamp       int64                   `bson:"timestamp"`
	Difficulty      string                  `bson:"difficulty"`
	TotalDifficulty string                  `bson:"totalDifficulty"`
	Creator         string                  `bson:"creator"`
	Nonce           string                  `bson:"nonce"`
	TxHash          string                  `bson:"txHash"`
	Reward          int64                   `bson:"reward"`
	Txs             []DBSimpleTxInBlock     `bson:"transactions"`
	Debt            []DBSimpleTxDebtInBlock `bson:"debt"`
	TxDebt          []DBSimpleTxDebtInBlock `bson:"txDebt"`
	ShardNumber     int                     `bson:"shardNumber"`
}

//Debt describle a transaction which stored in the database
type Debt struct {
	Hash        string `bson:"hash"`
	TxHash      string `bson:"txhash"`
	From        string `bson:"from"`
	To          string `bson:"to"`
	Height      uint64 `bson:"height"`
	Idx         uint64 `bson:"idx"`
	ShardNumber int    `bson:"shardNumber"`
	Fee         int64  `bson:"fee"`
	Payload     string `bson:"payload"`
	Amount      int64  `bson:"amount"`
}

type DBSimpleTxDebtInBlock struct {
	TxHash      string `bson:"txHash"`
	ShardNumber int    `bson:"shardNumber"`
	Account     string `bson:"account"`
	Amount      int64  `bson:"amount"`
	Fee         int64  `bson:"fee"`
	Payload     string `bson:"payload"`
}

//DBTx describle a transaction which stored in the database
type DBTx struct {
	TxType          int         `bson:"txtype"` // 0 is an normal transaction, 1 is an create contract transaction
	Hash            string      `bson:"hash"`
	DebtTxHash      string      `bson:"debtTxHash"`
	From            string      `bson:"from"`
	To              string      `bson:"to"`
	Amount          int64       `bson:"amount"`
	AccountNonce    string      `bson:"accountNonce"`
	Timestamp       string      `bson:"timestamp"`
	Timetxs         string      `bson:"timetxs"`
	Payload         string      `bson:"payload"`
	Block           uint64      `bson:"block"`
	Idx             int64       `bson:"idx"`
	ShardNumber     int         `bson:"shardNumber"`
	Fee             int64       `bson:"fee"`
	Pending         bool        `bson:"pending"`
	ContractAddress string      `bson:"contractAddress"`
	Receipt         rpc.Receipt `bson:"receipt"`
}

//DBAccount describle a account which stored in the database
type DBAccount struct {
	AccType     int    `bson:"accType"` //0 is normal account, 1 is contract account
	Address     string `bson:"address"`
	Balance     int64  `bson:"balance"`
	ShardNumber int    `bson:"shardNumber"`
	TxCount     int64  `bson:"txCount"`
	TimeStamp   int64  `bson:"timestamp"`
}

//DBMiner describle a miner account which stored in the database
type DBMiner struct {
	Address     string `bson:"address"`
	Revenue     int64  `bson:"total"`
	ShardNumber int    `bson:"shardNumber"`
	Reward      int64  `bson:"reward"`
	TxFee       int64  `bson:"fee"`
	TimeStamp   int64  `bson:"timestamp"`
}

//CreateDbBlock convert an rpc block to an dbblock
func CreateDbBlock(b *rpc.BlockInfo) *DBBlock {
	var dbBlock DBBlock
	dbBlock.HeadHash = b.Hash
	dbBlock.PreHash = b.ParentHash
	dbBlock.Height = int64(b.Height)
	dbBlock.Timestamp = b.Timestamp.Int64()
	dbBlock.Difficulty = b.Difficulty.String()
	dbBlock.TotalDifficulty = b.TotalDifficulty.String()
	dbBlock.Creator = b.Creator
	dbBlock.Nonce = strconv.FormatUint(b.Nonce, 10)
	//exclude coinbase transaction
	for i := 0; i < len(b.Txs); i++ {
		var simpleTx DBSimpleTxInBlock
		simpleTx.Hash = b.Txs[i].Hash
		simpleTx.From = b.Txs[i].From
		simpleTx.To = b.Txs[i].To
		simpleTx.Fee = b.Txs[i].Fee
		simpleTx.Amount = b.Txs[i].Amount.Int64()
		simpleTx.Timestamp = strconv.FormatUint(b.Txs[i].Timestamp, 10)
		dbBlock.Txs = append(dbBlock.Txs, simpleTx)

		if i != len(b.Txs)-1 {
			dbBlock.Reward += b.Txs[i].Fee
		}
	}

	//coinbase reward
	if len(b.Txs) > 0 {
		tx := b.Txs[len(b.Txs)-1]
		dbBlock.Reward = tx.Amount.Int64()
	}

	for i := 0; i < len(b.Debts); i++ {
		var simpleTxdebt DBSimpleTxDebtInBlock
		simpleTxdebt.Account = b.Debts[i].To
		simpleTxdebt.TxHash = b.Debts[i].TxHash
		simpleTxdebt.ShardNumber = b.Debts[i].ShardNumber
		simpleTxdebt.Fee = b.Debts[i].Fee
		simpleTxdebt.Payload = b.Debts[i].Payload
		simpleTxdebt.Amount = b.Debts[i].Amount.Int64()
		dbBlock.Debt = append(dbBlock.Debt, simpleTxdebt)
	}

	for i := 0; i < len(b.TxDebts); i++ {
		var simpleTxdebt DBSimpleTxDebtInBlock
		simpleTxdebt.Account = b.TxDebts[i].To
		simpleTxdebt.TxHash = b.TxDebts[i].TxHash
		simpleTxdebt.ShardNumber = b.TxDebts[i].ShardNumber
		simpleTxdebt.Fee = b.TxDebts[i].Fee
		simpleTxdebt.Payload = b.TxDebts[i].Payload
		simpleTxdebt.Amount = b.TxDebts[i].Amount.Int64()
		dbBlock.TxDebt = append(dbBlock.TxDebt, simpleTxdebt)
	}

	return &dbBlock
}

//CreateDbTx convert an rpc transaction to an dbtransaction
func CreateDbTx(t rpc.Transaction) *DBTx {
	var trans DBTx
	trans.TxType = t.TxType
	trans.Hash = t.Hash
	trans.DebtTxHash = t.DebtTxHash
	trans.From = t.From
	trans.To = t.To
	trans.Amount = t.Amount.Int64()
	timetxs := time.Unix(int64(t.Timestamp), 0)
	trans.Timetxs = timetxs.Format("2006-01-02")
	trans.Timestamp = strconv.FormatUint(t.Timestamp, 10)
	trans.AccountNonce = strconv.FormatUint(t.AccountNonce, 10)
	trans.Payload = t.Payload
	trans.Block = t.Block
	trans.Idx = int64(t.Idx)
	trans.Fee = t.Fee
	return &trans
}

//CreateDebtTx convert an rpc transaction to an dbtransaction
func CreateDebtTx(t rpc.Debt) *Debt {
	var debts Debt
	debts.Hash = t.Hash
	debts.TxHash = t.TxHash
	debts.To = t.To
	debts.Amount = t.Amount.Int64()
	debts.Payload = t.Payload
	debts.Height = t.Block
	debts.Fee = t.Fee
	return &debts
}

//CreateEmptyAccount create an empty dbaccount
func CreateEmptyAccount(address string, shardNumber int) *DBAccount {
	return &DBAccount{
		Address:     address,
		ShardNumber: shardNumber,
	}
}

//CreateDBAccountTx create an dbaccounttx stored in DBAccount for quickly search
/*
func CreateDBAccountTx(b *rpc.BlockInfo, rpcTx *rpc.Transaction, inOrOut bool) *DBAccountTx {
	tx := &DBAccountTx{
		Hash:      rpcTx.Hash,
		Block:     int64(b.Height),
		From:      rpcTx.From,
		To:        rpcTx.To,
		Amount:    rpcTx.Amount.Int64(),
		Timestamp: int64(rpcTx.Timestamp),
		TxFee:     0.0,
		InOrOut:   inOrOut,
	}

	return tx
}
*/

//DBOneDayTxInfo describle all transactions in an single day
type DBOneDayTxInfo struct {
	TotalTxs    int   `bson:"totaltxs"`
	TotalBlocks int   `bson:"totalblocks"`
	TimeStamp   int64 `bson:"timestamp"`
	ShardNumber int   `bson:"shardnumber"`
}

//DBOneDayHashRate describle all hashrates in an single day
type DBOneDayHashRate struct {
	HashRate    float64 `bson:"hashrate"`
	TimeStamp   int64   `bson:"timestamp"`
	ShardNumber int     `bson:"shardnumber"`
}

//DBOneDayBlockDifficulty describle avg block difficulty in an single day
type DBOneDayBlockDifficulty struct {
	Difficulty  float64 `bson:"difficulty"`
	TimeStamp   int64   `bson:"timestamp"`
	ShardNumber int     `bson:"shardnumber"`
}

//DBOneDayBlockAvgTime describle avg block time in an single day
type DBOneDayBlockAvgTime struct {
	AvgTime     float64 `bson:"avgtime"`
	TimeStamp   int64   `bson:"timestamp"`
	ShardNumber int     `bson:"shardnumber"`
}

//DBOneDayBlockInfo describle all blocks in an single day
type DBOneDayBlockInfo struct {
	TotalBlocks int64 `bson:"totalblocks"`
	Rewards     int64 `bson:"rewards"`
	TimeStamp   int64 `bson:"timestamp"`
	ShardNumber int   `bson:"shardnumber"`
}

//DBOneDayAddressInfo describle all blocks in an single day
type DBOneDayAddressInfo struct {
	TotalAddresss int64 `bson:"totaladdresss"`
	TodayIncrease int64 `bson:"todayincrease"`
	TimeStamp     int64 `bson:"timestamp"`
	ShardNumber   int   `bson:"shardnumber"`
}

//DBOneDaySingleAddressInfo describle one day single address info
type DBOneDaySingleAddressInfo struct {
	Address     string `bson:"address"`
	TimeStamp   int64  `bson:"timestamp"`
	ShardNumber int    `bson:"shardnumber"`
}

//DBSingleMinerRankInfo describle single miner rank info
type DBSingleMinerRankInfo struct {
	Address    string  `bson:"address"`
	Mined      int     `bson:"mined"`
	Percentage float64 `bson:"percentage"`
}

//DBMinerRankInfo descible top miner rank
type DBMinerRankInfo struct {
	Rank        []DBSingleMinerRankInfo `bson:"rank"`
	ShardNumber int                     `bson:"shardnumber"`
}

//DBNodeInfo descible an single node in the network
type DBNodeInfo struct {
	ShardNumber          int    `bson:"shardNumber"`
	ID                   string `bson:"id"`
	Host                 string `bson:"host"`
	Port                 string `bson:"port"`
	City                 string `bson:"city"`
	Region               string `bson:"region"`
	Country              string `bson:"country"`
	Client               string `bson:"client"`
	Caps                 string `bson:"caps"`
	LastSeen             int64  `bson:"lastseen"`
	LongitudeAndLatitude string `bson:"longitudeandlatitude"`
}
