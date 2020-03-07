package handlers

import "github.com/seeleteam/scan-api/database"

// BlockInfoDB Warpper for access mongodb.
type BlockInfoDB interface {
	GetBlockTxsTps() (float64, error)
	GetBlockHeight(shardNumber int) (uint64, error)
	GetBlockByHeight(shardNumber int, height uint64) (*database.DBBlock, error)
	GetTxsDayCount() ([]*database.DBTx, error)
	GetBlocksByHeight(shardNumber int, begin uint64, end uint64) ([]*database.DBBlock, error)
	GetBlockByHash(hash string) (*database.DBBlock, error)
	GetTxCnt() (uint64, error)
	GetBlockCnt() (uint64, error)
	GetBlockProTime() (int64, int64, error)
	GetAccountCnt() (uint64, error)
	GetContractCnt() (uint64, error)
	GetTxCntByShardNumber(shardNumber int) (uint64, error)
	GetdebtCntByShardNumber(shardNumber int) (uint64, error)
	GetPendingTxCntByShardNumber(shardNumber int) (uint64, error)
	GetTxByHash(hash string) (*database.DBTx, error)
	GetPendingTxByHash(hash string) (*database.DBTx, error)
	GetTxsByIdx(shardNumber int, begin uint64, end uint64) ([]*database.DBTx, error)
	GetdebtsByIdx(shardNumber int, begin uint64, end uint64) ([]*database.Debt, error)
	GetPendingTxsByIdx(shardNumber int, begin uint64, end uint64) ([]*database.DBTx, error)
	GetBlockfee(block uint64) (int64, error)
	GetTxsByAddresses(address string, asc bool, limit int, skip int) ([]*database.DBTx, error)
	GetPendingTxsByAddress(address string) ([]*database.DBTx, error)
	GetAccountCntByShardNumber(shardNumber int) (uint64, error)
	GetAccountByAddress(address string) (*database.DBAccount, error)
	GetAccountsByShardNumber(shardNumber int, max int) ([]*database.DBAccount, error)
	GetContractCntByShardNumber(shardNumber int) (uint64, error)
	GetContractsByShardNumber(shardNumber int, max int) ([]*database.DBAccount, error)
	GetTotalBalance() (map[int]int64, error)
	GetTxCntByShardNumberAndAddress(shardNumber int, address string) (int64, error)
	GetMinerAccounts(size int) ([]*database.DBMiner, error)
	GetAccountsByHome() []*database.DBAccount
	GetDebtByHash(hash string) (*database.Debt, error)
	GetblockdebtCntByShardNumber(shardNumber int, height uint64) (uint64, error)
	GetblockdebtsByIdx(shardNumber int, height uint64, begin uint64, end uint64) ([]*database.Debt, error)
	GetTxHis(startDate, today string) ([]*database.DBSimpleTxs, error)
	GetTxs(shardNumber int, sort string, desc bool , limit int, skip int) ([]*database.DBTx, error)
	UpdateContract(address string, sourceCode string, abiJson string) (error)
}

// ChartInfoDB Warpper for access mongodb.
type ChartInfoDB interface {
	GetTransInfoChart() ([]*database.DBOneDayTxInfo, error)
	GetOneDayAddressesChart() ([]*database.DBOneDayAddressInfo, error)
	GetOneDayBlockDifficultyChart() ([]*database.DBOneDayBlockDifficulty, error)
	GetOneDayBlocksChart() ([]*database.DBOneDayBlockInfo, error)
	GetHashRateChart() ([]*database.DBOneDayHashRate, error)
	GetOneDayBlockAvgTimeChart() ([]*database.DBOneDayBlockAvgTime, error)
	GetTopMinerChart() ([]*database.DBMinerRankInfo, error)

	GetTransInfoChartByShardNumber(shardNumber int) ([]*database.DBOneDayTxInfo, error)
	GetOneDayAddressesChartByShardNumber(shardNumber int) ([]*database.DBOneDayAddressInfo, error)
	GetOneDayBlockDifficultyChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockDifficulty, error)
	GetOneDayBlocksChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockInfo, error)
	GetHashRateChartByShardNumber(shardNumber int) ([]*database.DBOneDayHashRate, error)
	GetOneDayBlockAvgTimeChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockAvgTime, error)
	GetTopMinerChartByShardNumber(shardNumber int) ([]*database.DBMinerRankInfo, error)
}

// NodeInfoDB Warpper for access mongodb.
type NodeInfoDB interface {
	GetNodeInfosByShardNumber(shardNumber int) ([]*database.DBNodeInfo, error)
	GetNodeCntByShardNumber(shardNumber int) (uint64, error)
	GetNodeInfoByID(id string) (*database.DBNodeInfo, error)
}
