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
	"scan-api/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	blockItemNumsPrePage = 20
	transItemNumsPrePage = 25
	maxItemNumsPrePage   = 100

	blockTypestr = "block"
	transTypeStr = "transaction"
	accTypeStr   = "account"

	apiOk            = 0
	apiParmaInvalid  = 1
	apiInternalError = 2
	apiDBQueryError  = 3

	avgCountBlockNum = 5000
	txHashLength     = 66
	addressLength    = 130
)

var (
	errParamInvalid                     = errors.New("param is invalid")
	errGetBlockHeightFromDB             = errors.New("could not get block height from db")
	errGetTxCountFromDB                 = errors.New("could not get tx count from db")
	errGetBlockFromDB                   = errors.New("could not get block data from db")
	errGetTxFromDB                      = errors.New("could not get tx data from db")
	errGetAccountFromDB                 = errors.New("count not get account data from db")
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

//GetLastBlock get current block height
func GetLastBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		curBlockHeight, err := database.GetBlockHeight()
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		dbBlock, err := database.GetBlockByHeight(curBlockHeight - 1)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		age := getElpasedTimeDesc(big.NewInt(dbBlock.Timestamp))
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    age,
		})
	}
}

//GetBestBlock get current block height
func GetBestBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		curBlockHeight, err := database.GetBlockHeight()
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    curBlockHeight,
		})
	}
}

//GetHashRate get hash rate
func GetHashRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    database.Avg12HoursHashrate,
		})
	}
}

//GetDifficulty get difficulty
func GetDifficulty() gin.HandlerFunc {
	return func(c *gin.Context) {
		curBlockHeight, err := database.GetBlockHeight()
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		var dbBlock *database.DBBlock
		dbBlock, err = database.GetBlockByHeight(curBlockHeight - 1)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		ttDifficulty := big.NewInt(0)
		var avgDifficulty float64
		if ttDifficulty.UnmarshalText([]byte(dbBlock.TotalDifficulty)) == nil {
			avg := ttDifficulty.Div(ttDifficulty, big.NewInt(dbBlock.Height+1))
			avgDifficulty = float64(avg.Int64())
		} else {
			responseError(c, errDBDataError, http.StatusInternalServerError, apiInternalError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    avgDifficulty,
		})
	}
}

//GetAvgBlockTime get the latest 5000 blocks average time
func GetAvgBlockTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		curBlockHeight, err := database.GetBlockHeight()
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		curBlockHeight--
		var endBlock, beginBlock *database.DBBlock
		endBlock, err = database.GetBlockByHeight(curBlockHeight)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		var beginHeight uint64
		if curBlockHeight < avgCountBlockNum {
			beginHeight = 0
		} else {
			beginHeight = curBlockHeight - avgCountBlockNum
		}

		if beginHeight <= 0 {
			beginHeight = 1
		}
		beginBlock, err = database.GetBlockByHeight(beginHeight)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		timeElapsed := endBlock.Timestamp - beginBlock.Timestamp
		avgTime := (timeElapsed) / int64(curBlockHeight-beginHeight)
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    avgTime,
		})
	}
}

func getBlocksByBeginAndEnd(begin, end uint64) []*RetSimpleBlockInfo {
	var blocks []*RetSimpleBlockInfo
	dbBlocks, err := database.GetBlocksByHeight(begin, end)
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
func GetBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		if ps == 0 {
			ps = blockItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		curBlockHeight, err := database.GetBlockHeight()
		if err != nil {
			responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(curBlockHeight, p, ps)
		blocks := getBlocksByBeginAndEnd(begin, end)

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
func GetBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		height, err := strconv.ParseUint(c.Query("height"), 10, 64)
		if len(hash) > 0 {
			GetBlockDetailByHash(c, hash)
		} else if err == nil {
			GetBlockDetailByHeight(c, height)
		} else {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
		}
	}
}

//GetBlockDetailByHash get block by block hash
func GetBlockDetailByHash(c *gin.Context, hash string) {
	data, err := database.GetBlockByHash(hash)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := database.GetBlockHeight()

	detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data":    detailBlock,
	})
}

//GetBlockDetailByHeight get block by block height
func GetBlockDetailByHeight(c *gin.Context, height uint64) {
	data, err := database.GetBlockByHeight(height)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := database.GetBlockHeight()
	detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data":    detailBlock,
	})
}

//GetTxCnt handler for get all transaction count
func GetTxCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		transCnt, err := database.GetTxCnt()
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

//GetTxByHash handler for get transaction by hash
func GetTxByHash() gin.HandlerFunc {
	return func(c *gin.Context) {
		transHash := c.Query("txhash")
		if len(transHash) != txHashLength {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}
		data, err := database.GetTxByHash(transHash)
		if err != nil {
			responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			simpleTx := createRetSimpleTxInfo(data)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    simpleTx,
			})
		}

	}
}

func getTxsByBeginAndEnd(begin, end uint64) []*RetSimpleTxInfo {
	var txs []*RetSimpleTxInfo
	dbTrans, err := database.GetTxsByIdx(begin, end)
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
func GetTxsInBlock(c *gin.Context, height, p, ps uint64) {
	block, err := database.GetBlockByHeight(height)
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

//GetTxs get all transactions by order or by block
func GetTxs() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		if ps == 0 {
			ps = transItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		block, flag := c.GetQuery("c")
		if flag {
			height, err := strconv.ParseUint(block, 10, 64)
			if err != nil {
				responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			} else {
				GetTxsInBlock(c, height, p, ps)
				return
			}
		}

		txCnt, err := database.GetTxCnt()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(txCnt, p, ps)
		txs := getTxsByBeginAndEnd(begin, end)

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
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.Query("content")
		if content == "" {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}

		length := len(content)
		switch {
		case length == txHashLength:
			//found in transaction
			txHash := content
			data, err := database.GetTxByHash(txHash)
			if err != nil {
				responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
			} else {
				simpleTx := createRetSimpleTxInfo(data)

				c.JSON(http.StatusOK, gin.H{
					"code":    apiOk,
					"message": "",
					"data": gin.H{
						"type": transTypeStr,
						"info": simpleTx,
					},
				})
			}
		case length == addressLength:
			data, err := database.GetAccountByAddress(content)
			if err != nil {
				responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
			} else {
				detailAccount := createRetDetailAccountInfo(data)

				c.JSON(http.StatusOK, gin.H{
					"code":    apiOk,
					"message": "",
					"data": gin.H{
						"type": accTypeStr,
						"info": detailAccount,
					},
				})
			}
		default:
			//found in block
			blockNumber, err := strconv.ParseUint(content, 10, 64)
			if err != nil {
				responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
				return
			}

			data, err := database.GetBlockByHeight(blockNumber)
			if err != nil {
				responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			var maxHeight uint64
			maxHeight, err = database.GetBlockHeight()
			if err != nil {
				responseError(c, errGetBlockHeightFromDB, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data": gin.H{
					"type": blockTypestr,
					"info": detailBlock,
				},
			})
		}
	}
}

//getAccountsByBeginAndEnd
func getAccountsByBeginAndEnd(begin, end uint64) []*RetSimpleAccountInfo {
	var accounts []*RetSimpleAccountInfo
	dbAccounts := database.GetAccountsByIdx(begin, end)

	for i := 0; i < len(dbAccounts); i++ {
		data := dbAccounts[i]

		simpleAccount := createRetSimpleAccountInfo(data)
		simpleAccount.Rank = i + 1
		accounts = append(accounts, simpleAccount)
	}

	return accounts
}

//GetAccounts handler for get block list
func GetAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, _ := strconv.ParseUint(c.Query("p"), 10, 64)
		ps, _ := strconv.ParseUint(c.Query("ps"), 10, 64)
		if ps == 0 {
			ps = blockItemNumsPrePage
		} else if ps > maxItemNumsPrePage {
			ps = maxItemNumsPrePage
		}

		if p >= 1 {
			p--
		}

		accCnt := database.GetAccountCnt()

		page, begin, end := getBeginAndEndByPage(uint64(accCnt), p, ps)
		accounts := getAccountsByBeginAndEnd(begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": accCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": accounts,
			},
		})
	}
}

//GetAccountByAddress get account detail info by address
func GetAccountByAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		if len(address) != addressLength {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}

		data, err := database.GetAccountByAddress(address)
		if err != nil {
			responseError(c, errGetAccountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			detailAccount := createRetDetailAccountInfo(data)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    detailAccount,
			})
		}

	}
}
