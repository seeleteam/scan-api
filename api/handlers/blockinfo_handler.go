/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	blockItemNumsPrePage = 20
	transItemNumsPrePage = 25
	maxItemNumsPrePage   = 100

	blockTypestr    = "block"
	transTypeStr    = "transaction"
	accTypeStr      = "account"
	contractTypeStr = "contract"

	apiOk            = 0
	apiParmaInvalid  = 1
	apiInternalError = 2
	apiDBQueryError  = 3

	avgCountBlockNum = 5000
	txHashLength     = 66

	maxAccountTxCnt = 1000000
)

var (
	errParamInvalid                     = errors.New("param is invalid")
	errGetBlockHeightFromDB             = errors.New("could not get block height from db")
	errGetTxCountFromDB                 = errors.New("could not get tx count from db")
	errGetBlockFromDB                   = errors.New("could not get block data from db")
	errGetTxFromDB                      = errors.New("could not get tx data from db")
	errGetAccountFromDB                 = errors.New("count not get account data from db")
	errGetContractFromDB                = errors.New("count not get contract data from db")
	errDBDataError                      = errors.New("db data is error")
	errGetTxChartError                  = errors.New("could not get tx chart from db")
	errGetAddressChartError             = errors.New("could not get address chart from db")
	errGetBlockDifficultyChartError     = errors.New("could not get block difficulty chart from db")
	errGetBlockTimeChartError           = errors.New("could not get block time chart from db")
	errGetBlockCountAndRewardChartError = errors.New("could not get block count and reward chart from db")
	errGetHashRateChartError            = errors.New("could not get hashrate chart from db")
	errGetTopMinerChartError            = errors.New("could not get top miner chart from db")
	errGetNodeCountFromDB               = errors.New("could not get node count from db")
	errGetNodeInfoFromDB                = errors.New("could not get node data from db")
)

func responseError(c *gin.Context, err error, httpCode, code int) {
	fmt.Println(err)
	c.JSON(httpCode, gin.H{
		"code":    code,
		"message": err.Error(),
		"data":    gin.H{},
	})
}

//BlockHandler handle all block request
type BlockHandler struct {
	DBClient BlockInfoDB
}

//GetLastBlock get current block height
// func (h *BlockHandler) GetLastBlock() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		dbClinet := h.DBClient
// 		curBlockHeight, err := dbClinet.GetBlockHeight()
// 		if err != nil {
// 			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		dbBlock, err := dbClinet.GetBlockByHeight(curBlockHeight - 1)
// 		if err != nil {
// 			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		age := getElpasedTimeDesc(big.NewInt(dbBlock.Timestamp))
// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    apiOk,
// 			"message": "",
// 			"data":    age,
// 		})
// 	}
// }

//GetBestBlock get current block height
// func (h *BlockHandler) GetBestBlock() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		dbClinet := h.DBClient

// 		curBlockHeight, err := dbClinet.GetBlockHeight()
// 		if err != nil {
// 			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    apiOk,
// 			"message": "",
// 			"data":    curBlockHeight,
// 		})
// 	}
// }

//GetDifficulty get difficulty
// func (h *BlockHandler) GetDifficulty() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		dbClinet := h.DBClient

// 		curBlockHeight, err := dbClinet.GetBlockHeight()
// 		if err != nil {
// 			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		var dbBlock *database.DBBlock
// 		dbBlock, err = dbClinet.GetBlockByHeight(curBlockHeight - 1)
// 		if err != nil {
// 			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		ttDifficulty := big.NewInt(0)
// 		var avgDifficulty float64
// 		if ttDifficulty.UnmarshalText([]byte(dbBlock.TotalDifficulty)) == nil {
// 			avg := ttDifficulty.Div(ttDifficulty, big.NewInt(dbBlock.Height+1))
// 			avgDifficulty = float64(avg.Int64())
// 		} else {
// 			responseError(c, errDBDataError, http.StatusInternalServerError, apiInternalError)
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    apiOk,
// 			"message": "",
// 			"data":    avgDifficulty,
// 		})
// 	}
// }

//GetAvgBlockTime get the latest 5000 blocks average time
// func (h *BlockHandler) GetAvgBlockTime() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		dbClinet := h.DBClient

// 		curBlockHeight, err := dbClinet.GetBlockHeight()
// 		if err != nil {
// 			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		curBlockHeight--
// 		var endBlock, beginBlock *database.DBBlock
// 		endBlock, err = dbClinet.GetBlockByHeight(curBlockHeight)
// 		if err != nil {
// 			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		var beginHeight uint64
// 		if curBlockHeight < avgCountBlockNum {
// 			beginHeight = 0
// 		} else {
// 			beginHeight = curBlockHeight - avgCountBlockNum
// 		}

// 		if beginHeight <= 0 {
// 			beginHeight = 1
// 		}
// 		beginBlock, err = dbClinet.GetBlockByHeight(beginHeight)
// 		if err != nil {
// 			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
// 			return
// 		}

// 		timeElapsed := endBlock.Timestamp - beginBlock.Timestamp
// 		avgTime := (timeElapsed) / int64(curBlockHeight-beginHeight)
// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    apiOk,
// 			"message": "",
// 			"data":    avgTime,
// 		})
// 	}
// }

func (h *BlockHandler) getBlocksByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpleBlockInfo {
	dbClinet := h.DBClient

	var blocks []*RetSimpleBlockInfo
	dbBlocks, err := dbClinet.GetBlocksByHeight(shardNumber, begin, end)
	if err != nil {
		return nil
	}

	for i := 0; i < len(dbBlocks); i++ {
		data := dbBlocks[i]

		simpleBlock := createRetSimpleBlockInfo(data)
		blocks = append(blocks, simpleBlock)
	}

	return blocks
}

func getBeginAndEndByPage(total, p, step uint64) (page, begin, end uint64) {
	totalPages := uint64(math.Ceil(float64(total) / float64(step)))
	page = p
	if page > (totalPages - 1) {
		page = totalPages - 1
	}

	end = total - page*step
	if end < step {
		begin = 0
	} else {
		begin = end - step
	}

	return page, begin, end
}

//GetBlocks handler for get block list
func (h *BlockHandler) GetBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if ps == 0 {
			ps = blockItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		if s <= 0 {
			s = 1
		}

		shardNumber := int(s)
		curBlockHeight, err := dbClinet.GetBlockHeight(shardNumber)
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(curBlockHeight, p, ps)
		blocks := h.getBlocksByBeginAndEnd(shardNumber, begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": curBlockHeight,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": blocks,
			},
		})
	}
}

//GetBlock search block info by block height and block hash
func (h *BlockHandler) GetBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		hash, hashExist := c.GetQuery("hash")
		height, heightExist := c.GetQuery("height")
		s, _ := c.GetQuery("s")
		if hashExist {
			h.GetBlockDetailByHash(c, hash)
		} else if heightExist {
			blockheight, _ := strconv.ParseUint(height, 10, 64)
			s, _ := strconv.ParseInt(s, 10, 64)
			shaderNumber := int(s)
			h.GetBlockDetailByHeight(c, blockheight, shaderNumber)
		} else {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
		}
	}
}

//GetBlockDetailByHash get block by block hash
func (h *BlockHandler) GetBlockDetailByHash(c *gin.Context, hash string) {
	dbClinet := h.DBClient

	data, err := dbClinet.GetBlockByHash(hash)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := dbClinet.GetBlockHeight(data.ShardNumber)

	detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data":    detailBlock,
	})
}

//GetBlockDetailByHeight get block by block height
func (h *BlockHandler) GetBlockDetailByHeight(c *gin.Context, height uint64, shaderNumber int) {
	dbClinet := h.DBClient

	data, err := dbClinet.GetBlockByHeight(shaderNumber, height)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := dbClinet.GetBlockHeight(shaderNumber)
	detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data":    detailBlock,
	})
}

//GetTxCnt handler for get all transaction count
func (h *BlockHandler) GetTxCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		transCnt, err := dbClinet.GetTxCnt()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    transCnt,
			})
		}
	}
}

//GetBlockCnt handler for get all transaction count
func (h *BlockHandler) GetBlockCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		transCnt, err := dbClinet.GetBlockCnt()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    transCnt,
			})
		}
	}
}

//GetAccountCnt get all account count
func (h *BlockHandler) GetAccountCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		accountCnt, err := dbClinet.GetAccountCnt()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    accountCnt,
			})
		}
	}
}

//GetContractCnt get all contract count
func (h *BlockHandler) GetContractCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		contractCnt, err := dbClinet.GetContractCnt()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    contractCnt,
			})
		}
	}
}

//GetTxByHash handler for get transaction by hash
func (h *BlockHandler) GetTxByHash() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient
		transHash := c.Query("txhash")

		if len(transHash) != txHashLength {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}
		data, err := dbClinet.GetTxByHash(transHash)
		if err == nil {
			detailTx := createRetDetailTxInfo(data)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    detailTx,
			})
			return
		}

		data, err = dbClinet.GetPendingTxByHash(transHash)
		if err != nil {
			responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			simpleTx := createRetSimpleTxInfo(data)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    simpleTx,
			})
			return
		}

	}
}

func (h *BlockHandler) getTxsByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpleTxInfo {
	dbClinet := h.DBClient

	var txs []*RetSimpleTxInfo
	dbTrans, err := dbClinet.GetTxsByIdx(shardNumber, begin, end)
	if err != nil {
		return nil
	}
	for i := 0; i < len(dbTrans); i++ {
		data := dbTrans[i]

		simpleTransaction := createRetSimpleTxInfo(data)
		txs = append(txs, simpleTransaction)
	}

	return txs
}

func (h *BlockHandler) getPendingTxsByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpleTxInfo {
	dbClinet := h.DBClient

	var txs []*RetSimpleTxInfo
	dbTrans, err := dbClinet.GetPendingTxsByIdx(shardNumber, begin, end)
	if err != nil {
		return nil
	}
	for i := 0; i < len(dbTrans); i++ {
		data := dbTrans[i]

		simpleTransaction := createRetSimpleTxInfo(data)
		txs = append(txs, simpleTransaction)
	}

	return txs
}

//GetTxsInBlock get all transactions in block by height
func (h *BlockHandler) GetTxsInBlock(c *gin.Context, shardNumber int, height, p, ps uint64) {
	dbClinet := h.DBClient

	block, err := dbClinet.GetBlockByHeight(shardNumber, height)
	if err != nil {
		responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	txCntInBlock := len(block.Txs)
	page, begin, end := getBeginAndEndByPage(uint64(txCntInBlock), p, ps)
	txs := block.Txs[begin:end]

	var retTxs []*RetSimpleTxInfo
	for i := 0; i < len(txs); i++ {
		data := txs[i]

		timeStamp := big.NewInt(0)
		var age string
		if timeStamp.UnmarshalText([]byte(data.Timestamp)) == nil {
			age = getElpasedTimeDesc(timeStamp.Div(timeStamp, big.NewInt(1e9)))
		}

		simpleTransaction := &RetSimpleTxInfo{
			TxHash: data.Hash,
			Block:  height,
			Age:    age,
			From:   data.From,
			To:     data.To,
			Value:  data.Amount,
		}
		retTxs = append(retTxs, simpleTransaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data": gin.H{
			"pageInfo": gin.H{
				"totalCount": txCntInBlock,
				"begin":      begin,
				"end":        end,
				"curPage":    page + 1,
			},
			"list": retTxs,
		},
	})
}

//GetTxsInAccount get tx list from this account
func (h *BlockHandler) GetTxsInAccount(c *gin.Context, address string, p, ps uint64) {
	dbClinet := h.DBClient

	txs, err := dbClinet.GetTxsByAddresss(address, maxAccountTxCnt)
	if err != nil {
		responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	pengdingTxs, err := dbClinet.GetPendingTxsByAddress(address)
	if err != nil {
		responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	txs = append(pengdingTxs, txs...)

	txCntInAccount := len(txs)
	page, begin, end := getBeginAndEndByPage(uint64(txCntInAccount), p, ps)
	txs = txs[begin:end]

	var retTxs []*RetDetailAccountTxInfo
	for i := 0; i < len(txs); i++ {
		data := txs[i]

		timeStamp := big.NewInt(0)
		var age string
		if timeStamp.UnmarshalText([]byte(data.Timestamp)) == nil {
			age = getElpasedTimeDesc(timeStamp.Div(timeStamp, big.NewInt(1e9)))
		}

		var inOrOut bool
		if data.To == address {
			inOrOut = true
		} else {
			inOrOut = false
		}

		simpleTransaction := &RetDetailAccountTxInfo{
			ShardNumber: data.ShardNumber,
			TxType:      data.TxType,
			Hash:        data.Hash,
			Block:       data.Block,
			From:        data.From,
			To:          data.To,
			Value:       data.Amount,
			Age:         age,
			Fee:         data.Fee,
			InOrOut:     inOrOut,
			Pending:     data.Pending,
		}
		retTxs = append(retTxs, simpleTransaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data": gin.H{
			"pageInfo": gin.H{
				"totalCount": txCntInAccount,
				"begin":      begin,
				"end":        end,
				"curPage":    page + 1,
			},
			"list": retTxs,
		},
	})
}

//GetTxs get all transactions by order or by block
func (h *BlockHandler) GetTxs() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if ps == 0 {
			ps = transItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		if s <= 0 {
			s = 1
		}
		shardNumber := int(s)

		block, flag := c.GetQuery("block")
		if flag {
			height, err := strconv.ParseUint(block, 10, 64)
			if err != nil {
				responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			} else {
				h.GetTxsInBlock(c, shardNumber, height, p, ps)
				return
			}
		}

		address, flag := c.GetQuery("address")
		if flag {
			h.GetTxsInAccount(c, address, p, ps)
			return
		}

		txCnt, err := dbClinet.GetTxCntByShardNumber(shardNumber)
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(txCnt, p, ps)
		txs := h.getTxsByBeginAndEnd(shardNumber, begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": txCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": txs,
			},
		})
	}
}

//GetPendingTxs get an pending tx list
func (h *BlockHandler) GetPendingTxs() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if ps == 0 {
			ps = transItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		if s <= 0 {
			s = 1
		}
		shardNumber := int(s)

		txCnt, err := dbClinet.GetPendingTxCntByShardNumber(shardNumber)
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(txCnt, p, ps)
		txs := h.getPendingTxsByBeginAndEnd(shardNumber, begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": txCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": txs,
			},
		})
	}
}

//Search search something by transaction hash or block height
func (h *BlockHandler) Search(accHandler *AccountHandler, contractHanlder *ContractHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		content := c.Query("content")
		if content == "" {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}

		dbBlock, err := dbClinet.GetBlockByHash(content)
		if err == nil {
			var maxHeight uint64
			maxHeight, err = dbClinet.GetBlockHeight(dbBlock.ShardNumber)
			if err != nil {
				responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			detailBlock := createRetDetailBlockInfo(dbBlock, maxHeight, 0)
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": blockTypestr,
					"info": detailBlock,
				},
			})
			return
		}

		dbTx, err := dbClinet.GetTxByHash(content)
		if err == nil {
			detailTx := createRetDetailTxInfo(dbTx)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": transTypeStr,
					"info": detailTx,
				},
			})
			return
		}

		dbTx, err = dbClinet.GetPendingTxByHash(content)
		if err == nil {
			simpleTx := createRetSimpleTxInfo(dbTx)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": transTypeStr,
					"info": simpleTx,
				},
			})
			return
		}

		dbAccount := accHandler.GetAccountByAddressImpl(content)
		if dbAccount != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": accTypeStr,
					"info": dbAccount,
				},
			})
			return
		}

		dbContract := contractHanlder.GetContractByAddressImpl(content)
		if dbAccount != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": contractTypeStr,
					"info": dbContract,
				},
			})
			return
		}

		responseError(c, errParamInvalid, http.StatusOK, apiDBQueryError)
	}
}
