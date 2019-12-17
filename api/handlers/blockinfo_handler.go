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
	"strings"
	"time"

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
	errGetDebtFromDB                    = errors.New("could not get debt data from db")
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

func (h *BlockHandler) getBlocksByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpleBlockInfo {
	dbClient := h.DBClient

	var blocks []*RetSimpleBlockInfo
	dbBlocks, err := dbClient.GetBlocksByHeight(shardNumber, begin, end)
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

func getBeginAndEndByPageAndOrder(total, p, step uint64) (page, begin, end uint64) {
	totalPages := uint64(math.Ceil(float64(total) / float64(step)))
	page = p
	if page > (totalPages - 1) {
		page = totalPages - 1
	}

	end = (page + 1) * step
	if end >= total {
		end = total - 1
	}
	if end < step {
		begin = 0
	} else {
		begin = end - step
	}

	return page, begin, end
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

func getAccountBeginAndEndByPage(total, p, step uint64) (page, begin, end uint64) {
	totalPages := uint64(math.Ceil(float64(total) / float64(step)))
	page = p
	if page > (totalPages - 1) {
		page = totalPages - 1
	}

	end = total - page*step
	if end < step {
		begin = page * step
		end = total
	} else {
		end = (page + 1) * step
		begin = end - step
	}
	return page, begin, end
}

//GetBlocks handler for get block list
func (h *BlockHandler) GetBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

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
		curBlockHeight, err := dbClient.GetBlockHeight(shardNumber)
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
			shardNumber := int(s)
			h.GetBlockDetailByHeight(c, blockheight, shardNumber)
		} else {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
		}
	}
}

//GetBlockDetailByHash get block by block hash
func (h *BlockHandler) GetBlockDetailByHash(c *gin.Context, hash string) {
	dbClient := h.DBClient

	data, err := dbClient.GetBlockByHash(hash)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := dbClient.GetBlockHeight(data.ShardNumber)

	detailBlock := createRetDetailBlockInfo(data, maxHeight, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data":    detailBlock,
	})
}

//GetBlockDetailByHeight get block by block height
func (h *BlockHandler) GetBlockDetailByHeight(c *gin.Context, height uint64, shaderNumber int) {
	dbClient := h.DBClient

	data, err := dbClient.GetBlockByHeight(shaderNumber, height)
	if err != nil {
		responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	maxHeight, _ := dbClient.GetBlockHeight(shaderNumber)
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
		dbClient := h.DBClient

		transCnt, err := dbClient.GetTxCnt()
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
		dbClient := h.DBClient

		transCnt, err := dbClient.GetBlockCnt()
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

// GetBlockProTime get the last block information
func (h *BlockHandler) GetBlockProTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		lastBlockHeight, lastBlockTime, err := h.DBClient.GetBlockProTime()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			lastblockinfo := createRetLastblockInfo(lastBlockHeight, lastBlockTime)
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    lastblockinfo,
			})
		}
	}
}

//GetAccountCnt get all account count
func (h *BlockHandler) GetAccountCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

		accountCnt, err := dbClient.GetAccountCnt()
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
		dbClient := h.DBClient

		contractCnt, err := dbClient.GetContractCnt()
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

//GetBlockTxsTps TPS from block calculation
func (h *BlockHandler) GetBlockTxsTps() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.DBClient.GetBlockTxsTps()
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    data,
			})
		}
	}
}

//GetTxByHash handler for get transaction by hash
func (h *BlockHandler) GetTxByHash() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient
		transHash := c.Query("txhash")

		if len(transHash) != txHashLength {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}
		data, err := dbClient.GetTxByHash(transHash)
		if err == nil {
			detailTx := createRetDetailTxInfo(data)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    detailTx,
			})
			return
		}

		data, err = dbClient.GetPendingTxByHash(transHash)
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
	dbClient := h.DBClient

	var txs []*RetSimpleTxInfo
	dbTrans, err := dbClient.GetTxsByIdx(shardNumber, begin, end)
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

func (h *BlockHandler) getdebtsByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpledebtInfo {
	dbClient := h.DBClient

	var debts []*RetSimpledebtInfo
	debttxs, err := dbClient.GetdebtsByIdx(shardNumber, begin, end)
	if err != nil {
		return nil
	}

	for i := 0; i < len(debttxs); i++ {
		data := debttxs[i]

		simpleDebts := createRetSimpledebtInfo(data)
		debts = append(debts, simpleDebts)
	}

	return debts
}

func (h *BlockHandler) getblockdebtsByBeginAndEnd(shardNumber int, height uint64, begin, end uint64) []*RetSimpledebtInfo {
	dbClient := h.DBClient

	var debts []*RetSimpledebtInfo
	debt, err := dbClient.GetblockdebtsByIdx(shardNumber, height, begin, end)
	if err != nil {
		return nil
	}
	debtByBeginAndEnd := debt[begin:end]
	for i := 0; i < len(debtByBeginAndEnd); i++ {
		data := debtByBeginAndEnd[i]

		simpleDebts := createRetSimpledebtInfo(data)
		debts = append(debts, simpleDebts)
	}

	return debts
}

func (h *BlockHandler) getPendingTxsByBeginAndEnd(shardNumber int, begin, end uint64) []*RetSimpleTxInfo {
	dbClient := h.DBClient

	var txs []*RetSimpleTxInfo
	dbTrans, err := dbClient.GetPendingTxsByIdx(shardNumber, begin, end)
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

//GetTxsDayCount 30 days trading history data
func (h *BlockHandler) GetTxsDayCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
		todayDate := now.Format("2006-01-02")
		txs, err := h.DBClient.GetTxHis(startDate, todayDate)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		for _, tx := range txs {
			ymd := strings.Split(tx.Stime, "-")
			year, _ := strconv.Atoi(ymd[0])
			month, _ := strconv.Atoi(ymd[1])
			day, _ := strconv.Atoi(ymd[2])

			dateTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
			timestamp := dateTime.Unix()
			tx.Stime = strconv.FormatInt(timestamp, 10)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    txs,
		})

	}
}

//GetGasPrice Calculate the average gasprice values
func (h *BlockHandler) GetGasPrice() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		startDate := now.AddDate(0, 0, -10).Format("2006-01-02")
		todayDate := now.Format("2006-01-02")
		txs, err := h.DBClient.GetTxHis(startDate, todayDate)
		if err != nil {
			responseError(c, errGetBlockFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}
		var TxCount int
		var sumgas, highGasPrice, lowGasPrice int64

		highGasPrice = txs[0].HighGasPrice
		lowGasPrice = txs[0].LowGasPrice
		for i, tx := range txs {
			TxCount += tx.TxCount
			sumgas += tx.GasPrice
			if highGasPrice < txs[i].HighGasPrice {
				highGasPrice = txs[i].HighGasPrice
			}
			if lowGasPrice > txs[i].LowGasPrice {
				lowGasPrice = txs[i].LowGasPrice
			}

		}
		if TxCount == 0 {
			TxCount = 1
		}
		avegas := sumgas / int64(TxCount)
		walletgas := createwalletInfo(highGasPrice, lowGasPrice, avegas)
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    walletgas,
		})

	}
}

//GetTxsInBlock get all transactions in block by height
func (h *BlockHandler) GetTxsInBlock(c *gin.Context, shardNumber int, height, p, ps uint64) {
	dbClient := h.DBClient

	block, err := dbClient.GetBlockByHeight(shardNumber, height)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": 0,
					"begin":      0,
					"end":        0,
					"curPage":    0,
				},
				"list": nil,
			},
		})
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
			age = getElpasedTimeDesc(timeStamp)
		}

		simpleTransaction := &RetSimpleTxInfo{
			TxHash: data.Hash,
			Block:  height,
			Age:    age,
			From:   data.From,
			To:     data.To,
			Value:  data.Amount,
			Fee:	data.Fee,
			DebtHash:data.DebtTxHash,
			Timestamp:data.Timestamp,
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
func (h *BlockHandler) GetTxsInAccount(c *gin.Context, address string, p, ps int) {
	dbClient := h.DBClient
	txCntInAccount,err := dbClient.GetTxCntByShardNumberAndAddress(-1, address)
	txs, err := dbClient.GetTxsByAddresses(address,false,ps, p*ps)
	if err != nil {
		responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}
// no pending txs for now
/*	pengdingTxs, err := dbClient.GetPendingTxsByAddress(address)
	if err != nil {
		responseError(c, errGetTxFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}

	txs = append(pengdingTxs, txs...)*/

	var retTxs []*RetDetailAccountTxInfo
	for i := 0; i < len(txs); i++ {
		data := txs[i]

		timeStamp := big.NewInt(0)
		var age string
		if timeStamp.UnmarshalText([]byte(data.Timestamp)) == nil {
			age = getElpasedTimeDesc(timeStamp)
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
			Timestamp:	data.Timestamp,
		}
		retTxs = append(retTxs, simpleTransaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    apiOk,
		"message": "",
		"data": gin.H{
			"pageInfo": gin.H{
				"totalCount": txCntInAccount,
				"begin":      p*ps+1,
				"end":        (p+1)*ps,
				"curPage":    p + 1,
			},
			"list": retTxs,
		},
	})
}

//GetTxs get all transactions by order or by block
func (h *BlockHandler) GetTxs() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient
		p, _ := strconv.Atoi(c.Query("p"))
		ps, _ := strconv.Atoi(c.Query("ps"))
		s, _ := strconv.Atoi(c.Query("s"))
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

		// query transactions  for one block
		block, flag := c.GetQuery("block")
		if flag {
			height, err := strconv.ParseUint(block, 10, 64)
			if err != nil {
				responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			} else {
				h.GetTxsInBlock(c, shardNumber, height, uint64(p), uint64(ps))
				return
			}
		}
		//query transactions for one address
		address, flag := c.GetQuery("address")
		if flag {
			h.GetTxsInAccount(c, address, p, ps)
			return
		}
		//query transactions for one shard
		txCnt, err := dbClient.GetTxCntByShardNumber(shardNumber)
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}
		dbTrans,err := dbClient.GetTxs(shardNumber,"timestamp",true, ps,p*ps)
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}
		var txs []*RetSimpleTxInfo
		for i := 0; i < len(dbTrans); i++ {
			data := dbTrans[i]
			simpleTransaction := createRetSimpleTxInfo(data)
			txs = append(txs, simpleTransaction)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": txCnt,
					"begin":      p*ps+1,
					"end":        (p+1)*ps,
					"curPage":    p + 1,
				},
				"list": txs,
			},
		})
	}
}

//GetDebtByHash handler for get debt by hash
func (h *BlockHandler) GetDebtByHash() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient
		debtHash := c.Query("debtHash")

		if len(debtHash) != txHashLength {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}
		data, err := dbClient.GetDebtByHash(debtHash)
		if err != nil {
			responseError(c, errGetDebtFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		detailDebt := createRetDetailDebtInfo(data)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    detailDebt,
		})
		return
	}
}

//Getdebts get all debts by order or by block
func (h *BlockHandler) Getdebts() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

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

		debtCnt, err := dbClient.GetdebtCntByShardNumber(shardNumber)
		if err != nil {
			responseError(c, errGetTxCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(debtCnt, p, ps)
		debts := h.getdebtsByBeginAndEnd(shardNumber, begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": debtCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": debts,
			},
		})
	}
}

//GetBlockDebt get all debts in block by height
func (h *BlockHandler) GetBlockDebt() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

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
		height, _ := strconv.ParseUint(c.Query("block"), 10, 64)

		debtinblockCnt, err := dbClient.GetblockdebtCntByShardNumber(shardNumber, height)
		if err != nil {
			responseError(c, errGetDebtFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(debtinblockCnt, p, ps)
		debts := h.getblockdebtsByBeginAndEnd(shardNumber, height, begin, end)
		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": debtinblockCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": debts,
			},
		})
	}
}

//GetPendingTxs get an pending tx list
func (h *BlockHandler) GetPendingTxs() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

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

		txCnt, err := dbClient.GetPendingTxCntByShardNumber(shardNumber)
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
func (h *BlockHandler) Search(accHandler *AccountHandler, contractHandler *ContractHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClient := h.DBClient

		content := c.Query("content")
		if content == "" {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}

		dbBlock, err := dbClient.GetBlockByHash(content)
		if err == nil {
			var maxHeight uint64
			maxHeight, err = dbClient.GetBlockHeight(dbBlock.ShardNumber)
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

		dbTx, err := dbClient.GetTxByHash(content)
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

		dbTx, err = dbClient.GetPendingTxByHash(content)
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

		dbContract := contractHandler.GetContractByAddressImpl(content)
		if dbContract != nil {
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
