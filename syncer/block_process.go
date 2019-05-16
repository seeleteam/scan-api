package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
	"time"
)

func (s *Syncer) blockSync(block *rpc.BlockInfo) error {
	dbBlock := database.CreateDbBlock(block)
	var blockgas int64
	timeBegin := time.Now().Unix()
	for i := 0; i < len(dbBlock.Txs); i++ {
		trans := dbBlock.Txs[i]
		receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
		if err == nil {
			blockgas += receipt.UsedGas
			dbBlock.Txs[i].Fee = receipt.TotalFee
		}else{
			return err
		}
	}
	for i := 0; i < len(dbBlock.Debts); i++ {
		//receipt, err := s.rpc.GetReceiptByTxHash(dbBlock.Debts[i].TxHash)
	getDebt:	dbTx,err := s.db.GetTxByHash(dbBlock.Debts[i].TxHash)
		if err == nil {
			dbBlock.Debts[i].Fee = dbTx.Fee
		}else{
			time.Sleep(10*time.Second)
			log.Info("Try again to get debt's fee from transaction hash:%s",dbBlock.Debts[i].TxHash)
			goto getDebt
			return err
		}
	}
	log.Debug("seele_syncer block_process getReceiptHash time:%d(s)",time.Now().Unix()-timeBegin )
	dbBlock.UsedGas = blockgas
	dbBlock.ShardNumber = s.shardNumber
	// insert block info into database
	timeBegin = time.Now().Unix()
	if err := s.db.AddBlock(dbBlock); err != nil {
		return err
	}
	log.Debug("seele_syncer block_process addBlock to db time:%d(s)",time.Now().Unix()-timeBegin )
	// insert last block info into database to get final block produce rate
	timeBegin = time.Now().Unix()
	if err := storeLastBlocks(s.db, dbBlock); err != nil {
		return err
	}
	log.Debug("seele_syncer block_process storeLastBlocks time:%d(s)",time.Now().Unix()-timeBegin )
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
	if err != nil  || lastBlocks==nil{
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
