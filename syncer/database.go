package syncer

import "github.com/seeleteam/scan-api/database"

// Database wraps access to mongodb.
type Database interface {
	GetBlockHeight(shardNumber int) (uint64, error)
	AddBlock(b *database.DBBlock) error
	RemoveBlock(shard int, height uint64) error
	UpdateBlock(shard int, height uint64, b *database.DBBlock) error
	RemoveTxs(shard int, blockHeight uint64) error
	GetBlockByHeight(shardNumber int, height uint64) (*database.DBBlock, error)
	RemoveAllPendingTxs() error
	AddTx(tx *database.DBTx) error
	AddTxs(tx ...interface{}) error
	AddDebtTxs(debttxs ...interface{}) error
	AddPendingTx(tx *database.DBTx) error
	GetAccountByAddress(address string) (*database.DBAccount, error)
	GetMinerAccountByAddress(address string) (*database.DBMiner, error)
	AddAccount(account *database.DBAccount) error
	UpdateAccount(account *database.DBAccount) error
	UpdateMinerAccount(account *database.DBMiner) error
	UpdateAccountMinedBlock(address string, mined int64) error
	GetTxCntByShardNumber(shardNumber int) (uint64, error)
	GetPendingTxCntByShardNumber(shardNumber int) (uint64, error)
	GetTxCntByShardNumberAndAddress(shardNumber int, address string) (int64, error)
	GetMinedBlocksCntByShardNumberAndAddress(shardNumber int, address string) (int64, error)
}
