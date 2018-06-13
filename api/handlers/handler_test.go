/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
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

// type TestBlockInfoDB struct {
// }

// func (db *TestBlockInfoDB) GetBlockHeight(shardNumber int) (uint64, error) {
// 	return 1, nil
// }

// func (db *TestBlockInfoDB) GetBlockByHeight(shardNumber int, height uint64) (*database.DBBlock, error) {

// }

// func (db *TestBlockInfoDB) GetBlocksByHeight(shardNumber int, begin uint64, end uint64) ([]*database.DBBlock, error) {

// }

// func (db *TestBlockInfoDB) GetBlockByHash(hash string) (*database.DBBlock, error) {

// }

// func (db *TestBlockInfoDB) GetTxCnt() (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetBlockCnt() (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetAccountCnt() (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetContractCnt() (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetTxCntByShardNumber(shardNumber int) (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetPendingTxCntByShardNumber(shardNumber int) (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetTxByHash(hash string) (*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetPendingTxByHash(hash string) (*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetTxsByIdx(shardNumber int, begin uint64, end uint64) ([]*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetPendingTxsByIdx(shardNumber int, begin uint64, end uint64) ([]*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetTxsByAddresss(address string, max int) ([]*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetPendingTxsByAddress(address string) ([]*database.DBTx, error) {

// }

// func (db *TestBlockInfoDB) GetAccountCntByShardNumber(shardNumber int) (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetAccountByAddress(address string) (*database.DBAccount, error) {

// }

// func (db *TestBlockInfoDB) GetAccountsByShardNumber(shardNumber int, max int) ([]*database.DBAccount, error) {

// }

// func (db *TestBlockInfoDB) GetContractCntByShardNumber(shardNumber int) (uint64, error) {

// }

// func (db *TestBlockInfoDB) GetContractsByShardNumber(shardNumber int, max int) ([]*database.DBAccount, error) {

// }

// func (db *TestBlockInfoDB) GetTotalBalance() (map[int]int64, error) {

// }

// // ChartInfoDB Warpper for access mongodb.
// type TestChartInfoDB struct {
// }

// func (db *TestChartInfoDB) GetTransInfoChart() ([]*database.DBOneDayTxInfo, error) {

// }

// func (db *TestChartInfoDB) GetOneDayAddressesChart() ([]*database.DBOneDayAddressInfo, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlockDifficultyChart() ([]*database.DBOneDayBlockDifficulty, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlocksChart() ([]*database.DBOneDayBlockInfo, error) {

// }

// func (db *TestChartInfoDB) GetHashRateChart() ([]*database.DBOneDayHashRate, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlockAvgTimeChart() ([]*database.DBOneDayBlockAvgTime, error) {

// }

// func (db *TestChartInfoDB) GetTopMinerChart() ([]*database.DBMinerRankInfo, error) {

// }

// func (db *TestChartInfoDB) GetTransInfoChartByShardNumber(shardNumber int) ([]*database.DBOneDayTxInfo, error) {

// }

// func (db *TestChartInfoDB) GetOneDayAddressesChartByShardNumber(shardNumber int) ([]*database.DBOneDayAddressInfo, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlockDifficultyChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockDifficulty, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlocksChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockInfo, error) {

// }

// func (db *TestChartInfoDB) GetHashRateChartByShardNumber(shardNumber int) ([]*database.DBOneDayHashRate, error) {

// }

// func (db *TestChartInfoDB) GetOneDayBlockAvgTimeChartByShardNumber(shardNumber int) ([]*database.DBOneDayBlockAvgTime, error) {

// }

// func (db *TestChartInfoDB) GetTopMinerChartByShardNumber(shardNumber int) ([]*database.DBMinerRankInfo, error) {

// }

// type TestNodeInfoDB struct {
// }

// func (db *TestNodeInfoDB) GetNodeInfosByShardNumber(shardNumber int) ([]*database.DBNodeInfo, error) {

// }

// func (db *TestNodeInfoDB) GetNodeCntByShardNumber(shardNumber int) (uint64, error) {

// }

// func (db *TestNodeInfoDB) GetNodeInfoByID(id string) (*database.DBNodeInfo, error) {

// }

// func init() {

// }
