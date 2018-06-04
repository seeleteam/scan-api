/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"net/http"
	"scan-api/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

//getAccountsByBeginAndEnd
func getNodesByBeginAndEnd(begin, end uint64) []*database.DBNodeInfo {

	oneNodeInfos, _ := database.GetNodeInfos()

	if end > uint64(len(oneNodeInfos)) {
		end = uint64(len(oneNodeInfos)) - 1
	}

	nodes := oneNodeInfos[begin:end]

	return nodes
}

//GetNodes handler for get block list
func GetNodes() gin.HandlerFunc {
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

		nodeCnt, err := database.GetNodeCnt()
		if err != nil {
			responseError(c, errGetNodeCountFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		page, begin, end := getBeginAndEndByPage(uint64(nodeCnt), p, ps)
		nodes := getNodesByBeginAndEnd(begin, end)

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data": gin.H{
				"pageInfo": gin.H{
					"totalCount": nodeCnt,
					"begin":      begin,
					"end":        end,
					"curPage":    page + 1,
				},
				"list": nodes,
			},
		})
	}
}

//GetNode get node detail info by node id
func GetNode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		data, err := database.GetNodeInfoByID(id)
		if err != nil {
			responseError(c, errGetNodeInfoFromDB, http.StatusInternalServerError, apiDBQueryError)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    data,
			})
		}
	}
}

//GetNodeMap get node list
func GetNodeMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneNodeInfos, err := database.GetNodeInfos()
		if err != nil {
			responseError(c, errGetNodeInfoFromDB, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    oneNodeInfos,
		})
	}
}
