package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) blockSync(block *rpc.BlockInfo) error {
	log.Info("[BlockSync syncCnt:%d]Get Block %d", s.syncCnt, block.Height)

	//added block to cache
	dbBlock := database.CreateDbBlock(block)
	dbBlock.ShardNumber = s.shardNumber
	err := s.db.AddBlock(dbBlock)
	return err
}
