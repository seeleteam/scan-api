/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"scan-api/rpc"
	"strconv"
)

//DBSimpleTxInBlock describle the transaction info contained by dbblock which stored in the database
type DBSimpleTxInBlock struct {
	Hash      string `bson:"hash"`
	From      string `bson:"from"`
	To        string `bson:"to"`
	Amount    string `bson:"amount"`
	Timestamp string `bson:"timestamp"`
}

//DBBlock describle the block info which stored in the database
type DBBlock struct {
	HeadHash        string              `bson:"headHash"`
	PreHash         string              `bson:"preBlockHash"`
	Height          int64               `bson:"height"`
	StateHash       string              `bson:"stateHash"`
	Timestamp       int64               `bson:"timestamp"`
	Difficulty      string              `bson:"difficulty"`
	TotalDifficulty string              `bson:"totalDifficulty"`
	Creator         string              `bson:"creator"`
	Nonce           string              `bson:"nonce"`
	TxHash          string              `bson:"txHash"`
	Txs             []DBSimpleTxInBlock `bson:"transactions"`
}

//DBTx describle a transaction which stored in the database
type DBTx struct {
	Hash         string `bson:"hash"`
	From         string `bson:"from"`
	To           string `bson:"to"`
	Amount       string `bson:"amount"`
	AccountNonce string `bson:"accountNonce"`
	Timestamp    string `bson:"timestamp"`
	Payload      string `bson:"payload"`
	Block        string `bson:"block"`
	Idx          int64  `bson:"idx"`
}

//DBAccountTx describle a transaction contained by DBAccount
type DBAccountTx struct {
	Hash      string  `bson:"hash"`
	Block     int64   `bson:"block"`
	From      string  `bson:"from"`
	To        string  `bson:"to"`
	Amount    int64   `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
	TxFee     float64 `bson:"txfee"`
	InOrOut   bool    `bson:"inorout"`
}

//DBAccount describle a account which stored in the database
type DBAccount struct {
	Address string        `bson:"address"`
	Balance float64       `bson:"balance"`
	Txs     []DBAccountTx `bson:"txs"`
	TxCount int           `bson:"txcount"`
	Mined   int           `bson:"mined"`
}

//DBContract describle a contract which stored int the database
type DBContract struct {
	Address      string  `bson:"address"`
	ContractName string  `bson:"contractname"`
	Compiler     string  `bson:"compiler"`
	Balance      float64 `bson:"balance"`
	TxCount      int     `bson:"txCount"`
	DateVerified int64   `bson:"dateverified"`
}

//createDbBlock convert an rpc block to an dbblock
func createDbBlock(b *rpc.BlockInfo) *DBBlock {
	var dbBlock DBBlock
	dbBlock.HeadHash = b.Hash
	dbBlock.PreHash = b.ParentHash
	dbBlock.Height = int64(b.Height)
	dbBlock.Timestamp = b.Timestamp.Int64()
	dbBlock.Difficulty = b.Difficulty.String()
	dbBlock.TotalDifficulty = b.TotalDifficulty.String()
	dbBlock.Creator = b.Creator
	dbBlock.Nonce = strconv.FormatUint(b.Nonce, 10)
	for i := 0; i < len(b.Txs); i++ {
		var simpleTx DBSimpleTxInBlock
		simpleTx.Hash = b.Txs[i].Hash
		simpleTx.From = b.Txs[i].From
		simpleTx.To = b.Txs[i].To
		simpleTx.Amount = b.Txs[i].Amount.String()
		simpleTx.Timestamp = strconv.FormatUint(b.Txs[i].Timestamp, 10)
		dbBlock.Txs = append(dbBlock.Txs, simpleTx)
	}
	return &dbBlock
}

//createDbBlock convert an rpc transaction to an dbtransaction
func createDbTx(t rpc.Transaction) *DBTx {
	var trans DBTx
	trans.Hash = t.Hash
	trans.From = t.From
	trans.To = t.To
	trans.Amount = t.Amount.String()
	trans.Timestamp = strconv.FormatUint(t.Timestamp, 10)
	trans.AccountNonce = strconv.FormatUint(t.AccountNonce, 10)
	trans.Payload = t.Payload
	trans.Block = strconv.FormatUint(t.Block, 10)
	trans.Idx = int64(t.Idx)
	return &trans
}

//DBOneDayTxInfo describle all transactions in an single day
type DBOneDayTxInfo struct {
	TotalTxs    int   `bson:"totaltxs"`
	TotalBlocks int   `bson:"totalblocks"`
	TimeStamp   int64 `bson:"timestamp"`
}

//DBOneDayHashRate describle all hashrates in an single day
type DBOneDayHashRate struct {
	HashRate  float64 `bson:"hashrate"`
	TimeStamp int64   `bson:"timestamp"`
}

//DBOneDayBlockDifficulty describle avg block difficulty in an single day
type DBOneDayBlockDifficulty struct {
	Difficulty float64 `bson:"difficulty"`
	TimeStamp  int64   `bson:"timestamp"`
}

//DBOneDayBlockAvgTime describle avg block time in an single day
type DBOneDayBlockAvgTime struct {
	AvgTime   float64 `bson:"avgtime"`
	TimeStamp int64   `bson:"timestamp"`
}

//DBOneDayBlockInfo describle all blocks in an single day
type DBOneDayBlockInfo struct {
	TotalBlocks int64 `bson:"totalblocks"`
	Rewards     int64 `bson:"rewards"`
	TimeStamp   int64 `bson:"timestamp"`
}

//DBOneDayAddressInfo describle all blocks in an single day
type DBOneDayAddressInfo struct {
	TotalAddresss int64 `bson:"totaladdresss"`
	TodayIncrease int64 `bson:"todayincrease"`
	TimeStamp     int64 `bson:"timestamp"`
}

//DBOneDaySingleAddressInfo describle one day single address info
type DBOneDaySingleAddressInfo struct {
	Address   string `bson:"address"`
	TimeStamp int64  `bson:"timestamp"`
}

//DBSingleMinerRankInfo describle single miner rank info
type DBSingleMinerRankInfo struct {
	Address    string  `bson:"address"`
	Mined      int     `bson:"mined"`
	Percentage float64 `bson:"percentage"`
}

//DBMinerRankInfo descible top miner rank
type DBMinerRankInfo struct {
	Rank []DBSingleMinerRankInfo `bson:"rank"`
}

//DBNodeInfo descible an single node in the network
type DBNodeInfo struct {
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
