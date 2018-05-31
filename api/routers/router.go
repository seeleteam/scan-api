/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package routers

import (
	"scan-api/api/handlers"

	"github.com/gin-gonic/gin"
)

//InitRouters init all http handlers here
func InitRouters(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	v1.GET("/lastblock", handlers.GetLastBlock())
	v1.GET("/bestblock", handlers.GetBestBlock())
	v1.GET("/avgblocktime", handlers.GetAvgBlockTime())
	v1.GET("/block", handlers.GetBlock())
	v1.GET("/blocks", handlers.GetBlocks())
	v1.GET("/txcount", handlers.GetTxCnt())
	v1.GET("/txs", handlers.GetTxs())
	v1.GET("/tx", handlers.GetTxByHash())
	v1.GET("/search", handlers.Search())
	v1.GET("/accounts", handlers.GetAccounts())
	v1.GET("/account", handlers.GetAccountByAddress())

	v1.GET("/difficulty", handlers.GetDifficulty())
	v1.GET("/hashrate", handlers.GetHashRate())

	v1.GET("./nodes", handlers.GetNodes())
	v1.GET("./node", handlers.GetNode())
	v1.GET("./nodemap", handlers.GetNodeMap())

	chartGrp := v1.Group("/chart")
	chartGrp.GET("/tx", handlers.GetTxHistory())
	chartGrp.GET("/difficulty", handlers.GetEveryDayBlockDifficulty())
	chartGrp.GET("/address", handlers.GetEveryDayAddress())
	chartGrp.GET("/blocks", handlers.GetEveryDayBlock())
	chartGrp.GET("/hashrate", handlers.GetEveryHashRate())
	chartGrp.GET("/blocktime", handlers.GetEveryDayBlockTime())
	chartGrp.GET("/miner", handlers.GetTopMiners())
}
