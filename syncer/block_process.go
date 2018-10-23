package syncer

import (
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/rpc"
)

func (s *Syncer) blockSync(block *rpc.BlockInfo) error {
	dbBlock := database.CreateDbBlock(block)
	dbBlock.ShardNumber = s.shardNumber
	return s.db.AddBlock(dbBlock)
}
