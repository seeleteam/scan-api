/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/seeleteam/scan-api/database"

	"github.com/gin-gonic/gin"
)

//NodeHandler handle all block request
type NodeHandler struct {
	DBClient  NodeInfoDB
	nodeInfos [][]*database.DBNodeInfo
}

//NewNodeHandler new an nodehandler
func NewNodeHandler(nodeDB NodeInfoDB) *NodeHandler {
	nodeInfos := make([][]*database.DBNodeInfo, shardCount, shardCount)
	ret := &NodeHandler{
		DBClient:  nodeDB,
		nodeInfos: nodeInfos,
	}

	ret.UpdateImpl()

	return ret
}

//getAccountsByBeginAndEnd
func (h *NodeHandler) getNodesByBeginAndEnd(s int, begin, end uint64) []*database.DBNodeInfo {
	oneNodeInfos := h.nodeInfos[s-1]

	if end > uint64(len(oneNodeInfos)) {
		end = uint64(len(oneNodeInfos)) - 1
	}

	nodes := oneNodeInfos[begin:end]

	return nodes
}

//GetNodes handler for get block list
func (h *NodeHandler) GetNodes() gin.HandlerFunc {
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

		if s <= 0 || s > shardCount {
			responseError(c, errParamInvalid, http.StatusBadRequest, apiParmaInvalid)
		}

		nodeCnt := len(h.nodeInfos[s-1])

		page, begin, end := getBeginAndEndByPage(uint64(nodeCnt), p, ps)
		nodes := h.getNodesByBeginAndEnd(int(s), begin, end)

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
func (h *NodeHandler) GetNode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		for i := 0; i < shardCount; i++ {
			for j := 0; j < len(h.nodeInfos[i]); j++ {
				if h.nodeInfos[i][j].ID == id {
					c.JSON(http.StatusOK, gin.H{
						"code":    apiOk,
						"message": "",
						"data":    h.nodeInfos[i][j],
					})
					return
				}
			}
		}
		responseError(c, errGetNodeInfoFromDB, http.StatusInternalServerError, apiDBQueryError)
		return
	}
}

//GetNodeMap get node list
func (h *NodeHandler) GetNodeMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		var oneNodeInfos []*database.DBNodeInfo
		for i := 0; i < shardCount; i++ {
			oneNodeInfos = append(oneNodeInfos, h.nodeInfos[i]...)
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    oneNodeInfos,
		})
	}
}

//GetNodeCntChart get node count chart
func (h *NodeHandler) GetNodeCntChart() gin.HandlerFunc {
	return func(c *gin.Context) {

		nodeCntMap := make(map[int]int, len(h.nodeInfos))
		for i := 0; i < len(h.nodeInfos); i++ {
			nodeCntMap[i+1] = len(h.nodeInfos[i])
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    nodeCntMap,
		})
	}
}

//UpdateImpl cache nodes into memory
func (h *NodeHandler) UpdateImpl() {
	for i := 1; i <= shardCount; i++ {
		nodeInfos, err := h.DBClient.GetNodeInfosByShardNumber(i)
		h.nodeInfos[i-1] = make([]*database.DBNodeInfo, 0, 0)
		if err == nil {
			h.nodeInfos[i-1] = append(h.nodeInfos[i-1], nodeInfos...)
		}
	}
}

//Update set a timer to update node infos
func (h *NodeHandler) Update() {
	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Minute * 5)
		t := time.NewTimer(next.Sub(now))
		<-t.C

		h.UpdateImpl()
	}
}
