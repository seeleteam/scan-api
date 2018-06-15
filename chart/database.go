package chart

import "github.com/seeleteam/scan-api/database"

//ChartDB interface to access block db
type ChartDB interface {
	GetBlockByHeight(shardNumber int, height uint64) (*database.DBBlock, error)
	GetBlocksByTime(shardNumber int, beginTime, endTime int64) ([]*database.DBBlock, error)
	GetOneDayAddress(shardNumber int, zeroTime int64) (*database.DBOneDayAddressInfo, error)
	GetOneDaySingleAddressInfo(shardNumber int, address string) (*database.DBOneDaySingleAddressInfo, error)
	AddOneDaySingleAddressInfo(shardNumber int, t *database.DBOneDaySingleAddressInfo) error
	AddOneDayAddress(shardNumber int, t *database.DBOneDayAddressInfo) error
	GetOneDayBlock(shardNumber int, zeroTime int64) (*database.DBOneDayBlockInfo, error)
	AddOneDayBlock(shardNumber int, t *database.DBOneDayBlockInfo) error
	GetOneDayBlockDifficulty(shardNumber int, zeroTime int64) (*database.DBOneDayBlockDifficulty, error)
	AddOneDayBlockDifficulty(shardNumber int, t *database.DBOneDayBlockDifficulty) error
	GetOneDayBlockAvgTime(shardNumber int, zeroTime int64) (*database.DBOneDayBlockAvgTime, error)
	AddOneDayBlockAvgTime(shardNumber int, t *database.DBOneDayBlockAvgTime) error
	GetOneDayHashRate(shardNumber int, zeroTime int64) (*database.DBOneDayHashRate, error)
	AddOneDayHashRate(shardNumber int, t *database.DBOneDayHashRate) error
	RemoveTopMinerInfo() error
	AddTopMinerInfo(shardNumber int, rankInfo *database.DBMinerRankInfo) error
	AddOneDayTransInfo(shardNumber int, t *database.DBOneDayTxInfo) error
	GetOneDayTransInfo(shardNumber int, zeroTime int64) (*database.DBOneDayTxInfo, error)
}

var (
	GChartDB   ChartDB
	ShardCount int
)
