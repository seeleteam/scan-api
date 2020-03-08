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
	GasPrice  int64  `bson:"gasPrice"`
	GasLimit  int64  `bson:"gasLimit"`
	DebtTxHash string `bson:"debtTxHash"`
}

type DBSimpleTxs struct {
	Stime        string `json:"stime"`
	TxCount      int    `json:"txcount"`
	GasPrice     int64  `json:"gasprice"`
	HighGasPrice int64  `json:"highGasPrice"`
	LowGasPrice  int64  `json:"lowGasPrice"`
}

//DBBlock describle the block info which stored in the database
type DBBlock struct {
	HeadHash        string                `bson:"headHash"`
	PreHash         string                `bson:"preBlockHash"`
	Height          int64                 `bson:"height"`
	StateHash       string                `bson:"stateHash"`
	Timestamp       int64                 `bson:"timestamp"`
	Difficulty      string                `bson:"difficulty"`
	TotalDifficulty string                `bson:"totalDifficulty"`
	Creator         string                `bson:"creator"`
	Nonce           string                `bson:"nonce"`
	TxHash          string                `bson:"txHash"`
	Reward          int64                 `bson:"reward"`
	UsedGas         int64                 `bson:"usedGas"`
	Txs             []DBSimpleTxInBlock   `bson:"transactions"`
	Debts           []DBSimpleDebtInBlock `bson:"debt"`
	TxDebts         []DBSimpleDebtInBlock `bson:"txDebt"`
	ShardNumber     int                   `bson:"shardNumber"`
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

type DBSimpleDebtInBlock struct {
	Hash        string `bson:"hash"`
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
	GasPrice        int64       `bson:"gasPrice"`
	GasLimit        int64       `bson:"gasLimit"`
	UsedGas         int64       `bson:"usedGas"`
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
	SourceCode  string `bson:"sourceCode"`
	ABI         string `bson:"abi"`
}

//DBMiner describle a miner account which stored in the database
type DBMiner struct {
	Address     string `bson:"address"`
	Revenue     int64  `bson:"total"`
	ShardNumber int    `bson:"shardNumber"`
	Reward      int64  `bson:"reward"`
	TxFee       int64  `bson:"fee"`
	TimeStamp   int64  `bson:"timestamp"`
	Mined       int64  `bson:"mined"`
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
	dbBlock.StateHash = b.StateHash
	dbBlock.TxHash = b.TxHash
	//exclude coinbase transaction

	for i := 0; i < len(b.Txs); i++ {
		var simpleTx DBSimpleTxInBlock
		simpleTx.Hash = b.Txs[i].Hash
		simpleTx.From = b.Txs[i].From
		simpleTx.To = b.Txs[i].To
		simpleTx.Amount = b.Txs[i].Amount.Int64()
		simpleTx.Timestamp = b.Timestamp.String()
		simpleTx.GasLimit = b.Txs[i].GasLimit
		simpleTx.GasPrice = b.Txs[i].GasPrice
		dbBlock.Txs = append(dbBlock.Txs, simpleTx)

		//if i != len(b.Txs)-1 {
		//	dbBlock.Reward += b.Txs[i].Fee
		//}
		
		  if i != 0 {
			dbBlock.Reward += b.Txs[i].Fee
		}
		
	}

	//coinbase reward
	if len(b.Txs) > 0 {
		//tx := b.Txs[len(b.Txs)-1]
		tx := b.Txs[0]
		dbBlock.Reward = tx.Amount.Int64()
	}

	for i := 0; i < len(b.Debts); i++ {
		var simpleDebt DBSimpleDebtInBlock
		simpleDebt.Hash = b.Debts[i].Hash
		simpleDebt.Account = b.Debts[i].To
		simpleDebt.TxHash = b.Debts[i].TxHash
		simpleDebt.ShardNumber = b.Debts[i].ShardNumber
		simpleDebt.Fee = b.Debts[i].Fee
		simpleDebt.Payload = b.Debts[i].Payload
		simpleDebt.Amount = b.Debts[i].Amount.Int64()
		dbBlock.Debts = append(dbBlock.Debts, simpleDebt)
	}

	for i := 0; i < len(b.TxDebts); i++ {
		var simpleTxDebt DBSimpleDebtInBlock
		simpleTxDebt.Hash = b.TxDebts[i].Hash
		simpleTxDebt.Account = b.TxDebts[i].To
		simpleTxDebt.TxHash = b.TxDebts[i].TxHash
		simpleTxDebt.ShardNumber = b.TxDebts[i].ShardNumber
		simpleTxDebt.Fee = b.TxDebts[i].Fee
		simpleTxDebt.Payload = b.TxDebts[i].Payload
		simpleTxDebt.Amount = b.TxDebts[i].Amount.Int64()
		dbBlock.TxDebts = append(dbBlock.TxDebts, simpleTxDebt)
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
	timetxs := time.Unix(int64(t.Timestamp), 0).UTC()
	trans.Timetxs = timetxs.Format("2006-01-02")
	trans.Timestamp = strconv.FormatUint(t.Timestamp, 10)
	trans.AccountNonce = strconv.FormatUint(t.AccountNonce, 10)
	trans.Payload = t.Payload
	trans.Block = t.Block
	trans.Idx = int64(t.Idx)
	trans.Fee = t.Fee
	trans.GasPrice = t.GasPrice
	trans.GasLimit = t.GasLimit
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

// DBLastBlock contains the last block information
type DBLastBlock struct {
	ShardNumber int   `bson:"shardNumber"`
	Height      int64 `bson:"height"`
	Timestamp   int64 `bson:"timestamp"`
	TxNumber    int   `bson:"txNumber"`
}
