/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"scan-api/log"
	"scan-api/rpc"
	"time"
)

var (
	syncCnt uint64
)

//BlockSync get block data from seele node and store it in the mongodb
func BlockSync() error {
	log.Info("[BlockSync syncCnt:%d]Begin Sync", syncCnt)

	rpcSeeleRPC, err := rpc.GetSeeleRPC()
	if err != nil {
		rpc.ReleaseSeeleRPC()
		log.Error(err)
		return err
	}

	curBlock, err := rpcSeeleRPC.CurrentBlock()
	if err != nil {
		rpc.ReleaseSeeleRPC()
		log.Error(err)
		return err
	}

	curBlockHeight, err := GetBlockHeight()
	if err != nil {
		log.Error(err)
		return err
	}

	var flag bool
	if curBlockHeight < curBlock.Height {
		flag = true
	} else {
		flag = false
	}

	if flag {
		i := curBlockHeight
		for true {
			_, err := GetBlockByHeight(i)
			if err != nil {
				rpcSeeleRPC, err := rpc.GetSeeleRPC()
				if err != nil {
					log.Error(err)
					continue
				}

				rpcBlock, err := rpcSeeleRPC.GetBlockByHeight(i, true)
				if err != nil {
					rpc.ReleaseSeeleRPC()
					log.Error(err)
					continue
				}

				log.Info("[BlockSync syncCnt:%d]Get Block %d", syncCnt, i)

				//added block to cache
				err = AddBlock(createDbBlock(rpcBlock))
				if err != nil {
					log.Error(err)
					continue
				}

				for j := 0; j < len(rpcBlock.Txs); j++ {
					trans := rpcBlock.Txs[j]
					trans.Block = i
					transIdx, err := GetTxCnt()
					if err == nil {
						trans.Idx = transIdx
						AddTx(createDbTx(trans))
					} else {
						log.Error(err)
					}
				}

				ProcessAccount(rpcBlock)

				if i >= curBlock.Height {
					break
				} else {
					i++
				}
			} else {
				break
			}
		}

		ProcessGAccountTable()
	}

	if len(gAccountTbl) == 0 {
		ProcessGAccountTable()
	}

	ProcessLast12HoursHashRate()

	log.Info("[BlockSync syncCnt:%d]End Sync", syncCnt)
	return nil
}

//StartSync start an timer to sync block data from seele node
func StartSync(interval time.Duration) {
	ticks := time.NewTicker(interval * time.Second)
	tick := ticks.C
	go func() {
		for range tick {
			BlockSync()
			_, ok := <-tick
			if !ok {
				break
			}
		}
	}()
}
