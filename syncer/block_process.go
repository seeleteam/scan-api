package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) blockSync(block *rpc.BlockInfo) error {
	dbBlock := database.CreateDbBlock(block)
	//dbBlock.Txs[0]
	var blockgas int64
	for i := 0; i < len(dbBlock.Txs); i++ {
		trans := dbBlock.Txs[i]
		receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
		if err == nil {
			blockgas += receipt.UsedGas
			dbBlock.Txs[i].Fee = receipt.TotalFee
		}
	}

	dbBlock.UsedGas = blockgas
	dbBlock.ShardNumber = s.shardNumber
	// insert block info into database
	if err := s.db.AddBlock(dbBlock); err != nil {
		return err
	}
	// insert last block info into database to get final block produce rate
	if err := storeLastBlocks(s.db, dbBlock); err != nil {
		return err
	}
	return nil
}

func storeLastBlocks(db Database, block *database.DBBlock) error {
	// get last block
	last := &database.DBLastBlock{
		ShardNumber: block.ShardNumber,
		Height:      block.Height,
		Timestamp:   block.Timestamp,
		TxNumber:    len(block.Txs),
	}
	// get last two blocks by shard number
	lastBlocks, err := db.GetLastBlocksByShard(block.ShardNumber)
	if err != nil {
		// if blocks don't exist, insert the last block twice
		return initLastBlocks(db, last)
	}
	if len(lastBlocks) != 2 {
		// if the last blocks number is not 2, remove them all,
		// then insert the last block twice
		db.RemoveLastBlocksByShard(block.ShardNumber)
		return initLastBlocks(db, last)
	}
	// if the last block height is equal either the last two blocks, replace it,
	// or the lower height will be replaced by the last block
	var replace int64
	lastHeight := lastBlocks[0].Height
	secondHeight := lastBlocks[1].Height
	if lastHeight == last.Height || secondHeight == last.Height {
		replace = last.Height
	} else {
		replace = secondHeight
	}
	return db.UpdateLastBlock(replace, last)
}

func initLastBlocks(db Database, last *database.DBLastBlock) error {
	blocks := []interface{}{}
	blocks = append(blocks, last, last)
	return db.AddLastBlocks(blocks...)
}
