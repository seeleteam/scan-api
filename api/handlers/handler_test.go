/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"scan-api/log"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

type testResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Get According to the specific request uri, initiate a get request and return a response
func Get(uri string, router *gin.Engine) []byte {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// TestOnGetLastBlockRequest test the get last block handler
func TestOnGetLastBlockRequest(t *testing.T) {
	uri := "/api/v1/lastblock"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetBestBlockRequest(t *testing.T) {
	uri := "/api/v1/bestblock"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetAvgBlockTimeRequest(t *testing.T) {
	uri := "/api/v1/avgblocktimes"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetBlockRequest(t *testing.T) {
	uri := "/api/v1/block?height=1"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetBlocksRequest(t *testing.T) {
	uri := "/api/v1/blocks"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetTxCntRequest(t *testing.T) {
	uri := "/api/v1/txcount"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetTxsRequest(t *testing.T) {
	uri := "/api/v1/txs"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetTxByHashRequest(t *testing.T) {
	uri := "/api/v1/tx?txhash=0x60209b76ce6869a18266c8ba8608ac97394addd3bcdd3de3410cc5678c47d0b0"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Fatalf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnSearchRequest(t *testing.T) {
	uri := "/api/v1/search?content=0x60209b76ce6869a18266c8ba8608ac97394addd3bcdd3de3410cc5678c47d0b0"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetAccountsRequest(t *testing.T) {
	uri := "/api/v1/accounts"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetAccountByAddressRequest(t *testing.T) {
	uri := "/api/v1/account?address=0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetDifficultyRequest(t *testing.T) {
	uri := "/api/v1/difficulty"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetHashRateRequest(t *testing.T) {
	uri := "/api/v1/hashrate"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetNodesRequest(t *testing.T) {
	uri := "/api/v1/nodes"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetNodeRequest(t *testing.T) {
	uri := "/api/v1/node?id=99d08fe5c216335763277f26fdce972148f164ee5573afe1d59c50233754c53c083ba547767ebe71d82e291be5b6d077d17e20384fc8296430f71bb7b007f079"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetNodeMapRequest(t *testing.T) {
	uri := "/api/v1/nodemap"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartTxRequest(t *testing.T) {
	uri := "/api/v1/chart/tx"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartDifficultyRequest(t *testing.T) {
	uri := "/api/v1/chart/difficulty"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartAddressRequest(t *testing.T) {
	uri := "/api/v1/chart/address"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}
func TestOnGetChartBlockRequest(t *testing.T) {
	uri := "/api/v1/chart/blocks"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartHashRateRequest(t *testing.T) {
	uri := "/api/v1/chart/hashrate"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartBlockTimeRequest(t *testing.T) {
	uri := "/api/v1/chart/blocktime"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func TestOnGetChartTopMinersRequest(t *testing.T) {
	uri := "/api/v1/chart/miner"

	body := Get(uri, router)

	resp := new(testResponse)
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("respond error，body:%v\n", string(body))
	}

	if resp.Code != 0 {
		t.Errorf("respond error, error:%s\n", resp.Message)
	}
}

func init() {
	// 初始化路由
	router = gin.Default()

	log.NewLogger("debug", false)

	v1 := router.Group("/api/v1")
	v1.GET("/lastblock", GetLastBlock())
	v1.GET("/bestblock", GetBestBlock())
	v1.GET("/avgblocktime", GetAvgBlockTime())
	v1.GET("/block", GetBlock())
	v1.GET("/blocks", GetBlocks())
	v1.GET("/txcount", GetTxCnt())
	v1.GET("/txs", GetTxs())
	v1.GET("/tx", GetTxByHash())
	v1.GET("/search", Search())
	v1.GET("/accounts", GetAccounts())
	v1.GET("/account", GetAccountByAddress())

	v1.GET("/difficulty", GetDifficulty())
	v1.GET("/hashrate", GetHashRate())

	v1.GET("./nodes", GetNodes())
	v1.GET("./node", GetNode())
	v1.GET("./nodemap", GetNodeMap())

	chartGrp := v1.Group("/chart")
	chartGrp.GET("/tx", GetTxHistory())
	chartGrp.GET("/difficulty", GetEveryDayBlockDifficulty())
	chartGrp.GET("/address", GetEveryDayAddress())
	chartGrp.GET("/blocks", GetEveryDayBlock())
	chartGrp.GET("/hashrate", GetEveryHashRate())
	chartGrp.GET("/blocktime", GetEveryDayBlockTime())
	chartGrp.GET("/miner", GetTopMiners())
}
