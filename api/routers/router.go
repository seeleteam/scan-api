/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seeleteam/scan-api/api/handlers"
)

//Router api router
type Router struct {
	*handlers.AccountHandler
	*handlers.ContractHandler
	*handlers.BlockHandler
	*handlers.ChartHandler
	*handlers.NodeHandler
}

//New return an router
func New(blockDB handlers.BlockInfoDB, chartDB handlers.ChartInfoDB, nodeDB handlers.NodeInfoDB) *Router {
	accHandler := handlers.NewAccHandler(blockDB)
	contractHandler := handlers.NewContractHandler(blockDB)
	nodeHandler := handlers.NewNodeHandler(nodeDB)

	return &Router{
		AccountHandler:  accHandler,
		ContractHandler: contractHandler,
		BlockHandler:    &handlers.BlockHandler{DBClient: blockDB},
		ChartHandler:    &handlers.ChartHandler{DBClient: chartDB},
		NodeHandler:     nodeHandler,
	}
}

//Init init all http handlers here
func (r *Router) Init(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	//v1.GET("/lastblock", r.BlockHandler.GetLastBlock())
	//v1.GET("/bestblock", r.BlockHandler.GetBestBlock())
	//v1.GET("/avgblocktime", r.BlockHandler.GetAvgBlockTime())
	v1.GET("/accountcount", r.BlockHandler.GetAccountCnt())
	v1.GET("/block", r.BlockHandler.GetBlock())
	v1.GET("/blocks", r.BlockHandler.GetBlocks())
	v1.GET("/blockTxsTps", r.BlockHandler.GetBlockTxsTps())
	v1.GET("/blockprotime", r.BlockHandler.GetBlockProTime())
	v1.GET("/blockcount", r.BlockHandler.GetBlockCnt())
	v1.GET("/blockdebt", r.BlockHandler.GetBlockDebt())
	v1.GET("/contractcount", r.BlockHandler.GetContractCnt())
	v1.GET("/debts", r.BlockHandler.Getdebts())
	v1.GET("/debt", r.BlockHandler.GetDebtByHash())
	v1.GET("/Homeaccounts", r.AccountHandler.GetHomeAccounts())
	v1.GET("/pendingtxs", r.BlockHandler.GetPendingTxs())
	v1.GET("/txcount", r.BlockHandler.GetTxCnt())
	v1.GET("/txs", r.BlockHandler.GetTxs())
	v1.GET("/tx", r.BlockHandler.GetTxByHash())
	//ugly fix this
	v1.GET("/search", r.BlockHandler.Search(r.AccountHandler, r.ContractHandler))
	v1.GET("/accounts", r.AccountHandler.GetAccounts())
	v1.GET("/Txstat", r.BlockHandler.GetTxsDayCount())
	v1.GET("/account", r.AccountHandler.GetAccountByAddress())
	v1.GET("/miners", r.AccountHandler.GetMinerAccounts())
	v1.GET("/contracts", r.ContractHandler.GetContracts())
	v1.GET("/contract", r.ContractHandler.GetContractByAddress())
	v1.GET("/verifyContract", r.ContractHandler.VerifyContract())

	v1.GET("/Avegas", r.BlockHandler.GetGasPrice())
	//v1.GET("/difficulty", r.BlockHandler.GetDifficulty())
	//v1.GET("/hashrate", r.BlockHandler.GetHashRate())

	v1.GET("./nodes", r.NodeHandler.GetNodes())
	v1.GET("./node", r.NodeHandler.GetNode())
	v1.GET("./nodemap", r.NodeHandler.GetNodeMap())

	chartGrp := v1.Group("/chart")
	chartGrp.GET("/tx", r.ChartHandler.GetTxHistory())
	chartGrp.GET("/difficulty", r.ChartHandler.GetEveryDayBlockDifficulty())
	chartGrp.GET("/address", r.ChartHandler.GetEveryDayAddress())
	chartGrp.GET("/blocks", r.ChartHandler.GetEveryDayBlock())
	chartGrp.GET("/hashrate", r.ChartHandler.GetEveryHashRate())
	chartGrp.GET("/blocktime", r.ChartHandler.GetEveryDayBlockTime())
	chartGrp.GET("/miner", r.ChartHandler.GetTopMiners())
	chartGrp.GET("/node", r.NodeHandler.GetNodeCntChart())


	go r.AccountHandler.Update()
	go r.ContractHandler.Update()
	go r.NodeHandler.Update()
}
