package handlers

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
)

//ContractTbl describe
type ContractTbl struct {
	shardNumber   int
	DBClient      BlockInfoDB
	contractTbl   []*database.DBAccount
	contractMutex sync.RWMutex
	totalBalance  int64
}

//ProcessGContractTable process global account table
func (a *ContractTbl) ProcessGContractTable() {
	temp, err := a.DBClient.GetContractsByShardNumber(a.shardNumber, maxShowAccountNum)
	if err != nil {
		log.Error("[DB] err : %v", err)
	} else {
		a.contractMutex.Lock()
		a.contractTbl = temp
		a.contractMutex.Unlock()
	}
}

//GetContractCnt get the length of account table
func (a *ContractTbl) GetContractCnt() int {
	a.contractMutex.RLock()
	size := len(a.contractTbl)
	a.contractMutex.RUnlock()
	return size
}

//GetContractsByIdx get a transaction list from mongo by time period
func (a *ContractTbl) GetContractsByIdx(begin uint64, end uint64) []*database.DBAccount {
	a.contractMutex.RLock()
	if end > uint64(len(a.contractTbl)) {
		end = uint64(len(a.contractTbl)) - 1
	}

	retAccounts := a.contractTbl[begin:end]
	a.contractMutex.RUnlock()
	return retAccounts
}

func (a *ContractTbl) GetAccountByAddress(address string) *database.DBAccount {
	var retAccount *database.DBAccount
	a.contractMutex.RLock()
	for i := 0; i < len(a.contractTbl); i++ {
		if a.contractTbl[i].Address == address {
			retAccount = a.contractTbl[i]
		}
	}
	a.contractMutex.RUnlock()
	return retAccount
}

//getAccountsByBeginAndEnd
func (a *ContractTbl) getContractsByBeginAndEnd(begin, end uint64) []*RetSimpleAccountInfo {
	var accounts []*RetSimpleAccountInfo

	dbAccounts := a.GetContractsByIdx(begin, end)

	for i := 0; i < len(dbAccounts); i++ {
		data := dbAccounts[i]

		simpleAccount := createRetSimpleAccountInfo(data, a.totalBalance)
		simpleAccount.Rank = i + 1
		accounts = append(accounts, simpleAccount)
	}

	return accounts
}

//ContractHandler handle all contract request
type ContractHandler struct {
	contractTbls []*ContractTbl
	DBClient     BlockInfoDB
}

//NewContractHandler return an contractHandler to handler account request
func NewContractHandler(DBClient BlockInfoDB) *ContractHandler {
	var contractTbls []*ContractTbl
	for i := 1; i <= shardCount; i++ {
		contractTbl := &ContractTbl{shardNumber: i, DBClient: DBClient}
		contractTbls = append(contractTbls, contractTbl)
	}

	ret := &ContractHandler{
		contractTbls: contractTbls,
		DBClient:     DBClient,
	}
	ret.updateImpl()
	return ret
}

func (h *ContractHandler) updateImpl() {
	for i := 1; i <= shardCount; i++ {
		h.contractTbls[i-1].ProcessGContractTable()
	}

	totalBalances, err := h.DBClient.GetTotalBalance()
	if err != nil {
		log.Error("[DB] err : %v", err)
	}

	for i := 0; i < shardCount; i++ {
		if v, exist := totalBalances[i+1]; exist != false {
			h.contractTbls[i].totalBalance = v
		} else {
			h.contractTbls[i].totalBalance = remianTotalBalance
		}
	}
}

//Update Update account list every 5 secs
func (h *ContractHandler) Update() {
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Second * 5)
		t := time.NewTimer(next.Sub(now))
		<-t.C

		h.updateImpl()
	}
}

//GetContractByAddressImpl  use account info, account tx list and account pending tx list to assembly contract information
func (h *ContractHandler) GetContractByAddressImpl(address string) *RetDetailAccountInfo {
	dbClinet := h.DBClient

	data, err := dbClinet.GetAccountByAddress(address)
	if err != nil {
		return nil
	}

	if data.AccType != 1 {
		return nil
	}

	//txs, err := dbClinet.GetTxsByAddresss(address, txCount, false)
	txs, err := dbClinet.GetTxsByAddresses(address, false,txCount, -1)
	if err != nil {
		return nil
	}

	pengdingTxs, err := dbClinet.GetPendingTxsByAddress(address)
	if err != nil {
		return nil
	}

	var ttBalance int64
	if data.ShardNumber >= 1 && data.ShardNumber <= shardCount {
		ttBalance = h.contractTbls[data.ShardNumber-1].totalBalance
	}

	txs = append(txs, pengdingTxs...)

	data.TxCount, err = dbClinet.GetTxCntByShardNumberAndAddress(data.ShardNumber, address)
	if err != nil {
		return nil
	}

	if data.ShardNumber >= 1 && data.ShardNumber <= shardCount {
		account := h.contractTbls[data.ShardNumber-1].GetAccountByAddress(data.Address)
		if account != nil && account.TxCount != data.TxCount {
			account.TxCount = data.TxCount
		}
	}

	detailAccount := createRetDetailAccountInfo(data, txs, ttBalance)
	return detailAccount
}

//GetContractByAddress get contract detail info by address
func (h *ContractHandler) GetContractByAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		// if len(address) != addressLength {
		// 	responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
		// 	return
		// }

		detailAccount := h.GetContractByAddressImpl(address)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    detailAccount,
		})

	}
}

//GetContracts handler for get contract list
func (h *ContractHandler) GetContracts() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		if shardNumber < 1 || shardNumber > 20 {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
			return
		}

		contractTbl := h.contractTbls[shardNumber-1]
		contractCnt := contractTbl.GetContractCnt()

		page, begin, end := getBeginAndEndByPage(uint64(contractCnt), p, ps)
		contracts := contractTbl.getContractsByBeginAndEnd(begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": contractCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": contracts,
			},
		})
	}
}


func (h *ContractHandler) VerifyContract() gin.HandlerFunc{
	return func(c *gin.Context) {
		sourceCode := c.Query("sourceCode")
		abiJSON := c.Query("abi")
		address := c.Query("address")
		err := h.verifyContractImpl(address,sourceCode,abiJSON)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": err.Error(),
				"data":    false,
			})
		}else{
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    true,
			})
		}

	}
}

func (h *ContractHandler) verifyContractImpl(address string, sourceCode string,abiJSON string)(err error) {
	if address == "" || sourceCode == "" || abiJSON == "" {
		return errors.New("missing one or more query parameters")
	}
	err = h.DBClient.UpdateContract(address, sourceCode,abiJSON)
	if err !=nil {
		log.Error("save contract verification info failed, address:%s", address)
		return err
	}
	return nil
}