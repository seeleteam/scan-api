package syncer

import (
	"time"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) blockSync(block *rpc.BlockInfo) error {
	dbBlock := database.CreateDbBlock(block)
	var blockgas int64

	for i := 0; i < len(block.Txs); i++ {
		trans := block.Txs[i]
		receipt, err := s.rpc.GetReceiptByTxHash(trans.Hash)
		if err == nil {
			blockgas += receipt.UsedGas
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

func (s *Syncer) blockTxNumSync() {
	nTime := time.Now()
	startDate := nTime.AddDate(0, 0, -30).Format("2006-01-02")
	s.db.RemoveOutDateByDate(startDate)
	for i := 0; i < 30; i++ {
		yesTime := nTime.AddDate(0, 0, -i)
		yesTimeend := nTime.AddDate(0, 0, -i+1)
		logDay := yesTime.Format("20060102")
		logDayend := yesTimeend.Format("20060102")
		timeLayout := "20060102"
		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation(timeLayout, logDay, loc)
		theTimeend, _ := time.ParseInLocation(timeLayout, logDayend, loc)
		begin := theTime.Unix()
		end := theTimeend.Unix()

		dates := time.Unix(begin, 0)
		tx := new(database.DBSimpleTxs)
		tx.Stime = dates.Format("2006-01-02")

		cnt, err := s.db.GetTxsCntByDate(begin, end)
		if err != nil {
			cnt = 0
		}
		tx.TxCount = cnt
		s.db.UpdateTxsCntByDate(tx)
	}
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
