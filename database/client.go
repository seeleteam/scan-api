/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"errors"
	"scan-api/log"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	blockTbl = "block"
	txTbl    = "transaction"
	accTbl   = "account"

	chartTxTbl              = "chart_transhistory"
	chartHashRateTbl        = "chart_hashrate"
	chartBlockDifficultyTbl = "chart_blockdifficulty"
	chartBlockAvgTimeTbl    = "chart_blockavgtime"
	chartBlockTbl           = "chart_block"
	chartAddressTbl         = "chart_address"
	chartSingleAddressTbl   = "chart_single_address"
	chartTopMinerRankTbl    = "chart_topminer"

	nodeInfoTbl = "nodeinfo"
)

var (
	mgoSession *mgo.Session
	//DataBaseName mongo database name
	DataBaseName = "seele"
	//ConnURL mongo database address
	ConnURL = "127.0.0.1:27017"
	//db connect error
	errDBConnect = errors.New("could not connect to database")
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(ConnURL)
		if err != nil {
			log.Error("[DB] err : %v", err)
			return nil
		}
	}

	return mgoSession.Clone()
}

//InitDB init database connection
func InitDB() bool {
	mgo := getSession()
	if mgo != nil {
		return true
	}
	return false
}

//withCollection perform an database query
func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer func() {
		if session != nil {
			session.Close()
		}
	}()
	if session != nil {
		c := session.DB(DataBaseName).C(collection)
		err := s(c)
		processDataBaseError(err)
		return err
	}
	log.Error("[DB] err : could not connect to db, host is %s", ConnURL)
	return errDBConnect
}

//dropCollection test use remove the tbl
func dropCollection(tbl string) error {
	session := getSession()
	defer func() {
		if session != nil {
			session.Close()
		}
	}()
	if session != nil {
		c := session.DB(DataBaseName).C(tbl)
		err := c.DropCollection()
		processDataBaseError(err)
		return err
	}
	log.Error("[DB] err : could not connect to db, host is %s", ConnURL)
	return errDBConnect
}

//AddBlock insert a block into database
func AddBlock(b *DBBlock) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(b)
	}
	err := withCollection(blockTbl, query)
	return err
}

//removeBlock test use  remove block by height from database
func removeBlock(height uint64) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"height": height})
	}
	err := withCollection(blockTbl, query)
	return err
}

//GetBlockByHeight get block from mongo by block height
func GetBlockByHeight(height uint64) (*DBBlock, error) {
	b := new(DBBlock)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"height": height}).One(b)
	}
	err := withCollection(blockTbl, query)
	return b, err
}

//GetBlockByHash get a block from mongo by block header hash
func GetBlockByHash(hash string) (*DBBlock, error) {
	b := new(DBBlock)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"headHash": hash}).One(b)
	}
	err := withCollection(blockTbl, query)
	return b, err
}

//GetBlocksByHeight get a block list from mongo by height range
func GetBlocksByHeight(begin uint64, end uint64) ([]*DBBlock, error) {
	var blocks []*DBBlock

	query := func(c *mgo.Collection) error {

		return c.Find(bson.M{"height": bson.M{"$gte": begin, "$lt": end}}).Sort("-height").All(&blocks)
	}
	err := withCollection(blockTbl, query)
	return blocks, err
}

//GetBlocksByTime get a block list from mongo by time period
func GetBlocksByTime(beginTime, endTime int64) ([]*DBBlock, error) {
	var blocks []*DBBlock

	query := func(c *mgo.Collection) error {

		return c.Find(bson.M{"timestamp": bson.M{"$gte": beginTime, "$lte": endTime}}).All(&blocks)
	}
	err := withCollection(blockTbl, query)
	return blocks, err
}

//GetBlockHeight get row count of block table from mongo
func GetBlockHeight() (uint64, error) {
	var blockCnt uint64
	query := func(c *mgo.Collection) error {
		var err error
		//TODO: fix this overflow
		var temp int
		temp, err = c.Count()
		blockCnt = uint64(temp)
		return err
	}
	err := withCollection(blockTbl, query)
	return blockCnt, err
}

//AddTx insert a transaction into mongo
func AddTx(tx *DBTx) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(tx)
	}
	err := withCollection(txTbl, query)
	return err
}

//removeTx test use  remove tx by index from database
func removeTx(idx uint64) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"idx": strconv.FormatUint(idx, 10)})
	}
	err := withCollection(txTbl, query)
	return err
}

//GetTxByIdx get transaction from mongo by idx
func GetTxByIdx(idx uint64) (*DBTx, error) {
	tx := new(DBTx)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"idx": idx}).One(tx)
	}
	err := withCollection(txTbl, query)
	return tx, err
}

//GetTxsByIdx get a transaction list from mongo by time period
func GetTxsByIdx(begin uint64, end uint64) ([]*DBTx, error) {
	var trans []*DBTx
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"idx": bson.M{"$gte": begin, "$lt": end}}).Sort("-idx").All(&trans)
	}
	err := withCollection(txTbl, query)
	return trans, err
}

//GetTxByHash get transaction info by hash from mongo
func GetTxByHash(hash string) (*DBTx, error) {
	tx := new(DBTx)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"hash": hash}).One(tx)
	}
	err := withCollection(txTbl, query)
	return tx, err
}

//GetTxCnt get row count of transaction table from mongo
func GetTxCnt() (uint64, error) {
	var txCnt uint64
	query := func(c *mgo.Collection) error {
		var err error
		//TODO: fix this overflow
		var temp int
		temp, err = c.Count()
		txCnt = uint64(temp)
		return err
	}
	err := withCollection(txTbl, query)
	return txCnt, err
}

//removeAccount test use  remove account by address from database
func removeAccount(address string) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"address": address})
	}
	err := withCollection(accTbl, query)
	return err
}

//GetAccountByAddress get an dbaccount by account address
func GetAccountByAddress(address string) (*DBAccount, error) {
	account := new(DBAccount)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"address": address}).One(account)
	}
	err := withCollection(accTbl, query)
	return account, err
}

//AddAccount insert an account into database
func AddAccount(account *DBAccount) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(account)
	}
	err := withCollection(accTbl, query)
	return err
}

//UpdateAccount update account
func UpdateAccount(address string, balance float64, txs *[]DBAccountTx) error {
	incBalanceQuery := func(c *mgo.Collection) error {
		return c.Update(bson.M{"address": address},
			bson.M{"$inc": bson.M{
				"balance": balance,
			}})
	}
	err := withCollection(accTbl, incBalanceQuery)
	if err != nil {
		return err
	}

	addTxCountQuery := func(c *mgo.Collection) error {
		return c.Update(bson.M{"address": address},
			bson.M{"$inc": bson.M{
				"txcount": 1,
			}})
	}
	err = withCollection(accTbl, addTxCountQuery)
	if err != nil {
		return err
	}

	setTxQuery := func(c *mgo.Collection) error {
		return c.Update(bson.M{"address": address},
			bson.M{"$set": bson.M{
				"txs": txs,
			}})
	}

	err = withCollection(accTbl, setTxQuery)
	return err
}

//UpdateAccountMinedBlock update field mined block in the account info
func UpdateAccountMinedBlock(address string, mined int) error {
	incBalanceQuery := func(c *mgo.Collection) error {
		return c.Update(bson.M{"address": address},
			bson.M{"$inc": bson.M{
				"mined": mined,
			}})
	}
	err := withCollection(accTbl, incBalanceQuery)
	return err
}

//GetAccounts get an dbaccount list sort by balance
func GetAccounts(max int) ([]*DBAccount, error) {
	var accounts []*DBAccount
	query := func(c *mgo.Collection) error {
		return c.Find(nil).Sort("-balance").Limit(max).All(&accounts)
	}
	err := withCollection(accTbl, query)
	return accounts, err
}

//processDataBaseError shutdown database connection and log it
func processDataBaseError(err error) {
	if err == nil || err == mgo.ErrNotFound || err == mgo.ErrCursor {
		return
	}

	log.Error("[DB] err : %v", err)
	mgoSession.Close()
	mgoSession = nil
}

//AddOneDayTransInfo insert one dya transaction info into mongo
func AddOneDayTransInfo(t *DBOneDayTxInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartTxTbl, query)
	return err
}

//GetOneDayTransInfo get one day transaction info from mongo by zero hour timestamp
func GetOneDayTransInfo(zeroTime int64) (*DBOneDayTxInfo, error) {
	oneDayTransInfo := new(DBOneDayTxInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayTransInfo)
	}
	err := withCollection(chartTxTbl, query)
	return oneDayTransInfo, err
}

//GetTransInfoChart get all rows int the transhistory table
func GetTransInfoChart() ([]*DBOneDayTxInfo, error) {
	var oneDayTrans []*DBOneDayTxInfo
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayTrans)
	}
	err := withCollection(chartTxTbl, query)
	return oneDayTrans, err
}

//AddOneDayHashRate insert one dya hashrate info into mongo
func AddOneDayHashRate(t *DBOneDayHashRate) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartHashRateTbl, query)
	return err
}

//GetOneDayHashRate get one day hashrate info from mongo by zero hour timestamp
func GetOneDayHashRate(zeroTime int64) (*DBOneDayHashRate, error) {
	oneDayHashRate := new(DBOneDayHashRate)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayHashRate)
	}
	err := withCollection(chartHashRateTbl, query)
	return oneDayHashRate, err
}

//GetHashRateChart get all rows int the hashrate table
func GetHashRateChart() ([]*DBOneDayHashRate, error) {
	var oneDayHashRates []*DBOneDayHashRate
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayHashRates)
	}
	err := withCollection(chartHashRateTbl, query)
	return oneDayHashRates, err
}

//AddOneDayBlockDifficulty insert one dya avg block difficulty info into mongo
func AddOneDayBlockDifficulty(t *DBOneDayBlockDifficulty) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartBlockDifficultyTbl, query)
	return err
}

//GetOneDayBlockDifficulty get one day hashrate info from mongo by zero hour timestamp
func GetOneDayBlockDifficulty(zeroTime int64) (*DBOneDayBlockDifficulty, error) {
	oneDayBlockDifficulty := new(DBOneDayBlockDifficulty)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayBlockDifficulty)
	}
	err := withCollection(chartBlockDifficultyTbl, query)
	return oneDayBlockDifficulty, err
}

//GetOneDayBlockDifficultyChart get all rows int the hashrate table
func GetOneDayBlockDifficultyChart() ([]*DBOneDayBlockDifficulty, error) {
	var oneDayBlockDifficulties []*DBOneDayBlockDifficulty
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayBlockDifficulties)
	}
	err := withCollection(chartBlockDifficultyTbl, query)
	return oneDayBlockDifficulties, err
}

//AddOneDayBlockAvgTime insert one dya avg block time info into mongo
func AddOneDayBlockAvgTime(t *DBOneDayBlockAvgTime) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartBlockAvgTimeTbl, query)
	return err
}

//GetOneDayBlockAvgTime get one day avg block time info from mongo by zero hour timestamp
func GetOneDayBlockAvgTime(zeroTime int64) (*DBOneDayBlockAvgTime, error) {
	oneDayBlockAvgTime := new(DBOneDayBlockAvgTime)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayBlockAvgTime)
	}
	err := withCollection(chartBlockAvgTimeTbl, query)
	return oneDayBlockAvgTime, err
}

//GetOneDayBlockAvgTimeChart get all rows int the hashrate table
func GetOneDayBlockAvgTimeChart() ([]*DBOneDayBlockAvgTime, error) {
	var oneDayBlockAvgTimes []*DBOneDayBlockAvgTime
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayBlockAvgTimes)
	}
	err := withCollection(chartBlockAvgTimeTbl, query)
	return oneDayBlockAvgTimes, err
}

//AddOneDayBlock insert one dya block info into mongo
func AddOneDayBlock(t *DBOneDayBlockInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartBlockTbl, query)
	return err
}

//GetOneDayBlock get one day block info from mongo by zero hour timestamp
func GetOneDayBlock(zeroTime int64) (*DBOneDayBlockInfo, error) {
	oneDayBlock := new(DBOneDayBlockInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayBlock)
	}
	err := withCollection(chartBlockTbl, query)
	return oneDayBlock, err
}

//GetOneDayBlocksChart get all rows int the hashrate table
func GetOneDayBlocksChart() ([]*DBOneDayBlockInfo, error) {
	var oneDayBlocks []*DBOneDayBlockInfo
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayBlocks)
	}
	err := withCollection(chartBlockTbl, query)
	return oneDayBlocks, err
}

//AddOneDayAddress insert one dya block info into mongo
func AddOneDayAddress(t *DBOneDayAddressInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartAddressTbl, query)
	return err
}

//GetOneDayAddress get one day block info from mongo by zero hour timestamp
func GetOneDayAddress(zeroTime int64) (*DBOneDayAddressInfo, error) {
	oneDayAddress := new(DBOneDayAddressInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"timestamp": zeroTime}).One(oneDayAddress)
	}
	err := withCollection(chartAddressTbl, query)
	return oneDayAddress, err
}

//GetOneDayAddressesChart get all rows int the address table
func GetOneDayAddressesChart() ([]*DBOneDayAddressInfo, error) {
	var oneDayAddresses []*DBOneDayAddressInfo
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("timestamp").All(&oneDayAddresses)
	}
	err := withCollection(chartAddressTbl, query)
	return oneDayAddresses, err
}

//AddOneDaySingleAddressInfo insert one dya single address info into mongo
func AddOneDaySingleAddressInfo(t *DBOneDaySingleAddressInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err := withCollection(chartSingleAddressTbl, query)
	return err
}

//GetOneDaySingleAddressInfo get one day block info from mongo by zero hour timestamp
func GetOneDaySingleAddressInfo(address string) (*DBOneDaySingleAddressInfo, error) {
	oneDaySingleAddress := new(DBOneDaySingleAddressInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"address": address}).One(oneDaySingleAddress)
	}
	err := withCollection(chartSingleAddressTbl, query)
	return oneDaySingleAddress, err
}

//RemoveTopMinerInfo remove last 7 days top miner info
func RemoveTopMinerInfo() error {
	query := func(c *mgo.Collection) error {
		return c.DropCollection()
	}
	err := withCollection(chartTopMinerRankTbl, query)
	return err
}

//AddTopMinerInfo add top miner rank info into database
func AddTopMinerInfo(rankInfo *DBMinerRankInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(rankInfo)
	}
	err := withCollection(chartTopMinerRankTbl, query)
	return err
}

//GetTopMinerChart get all rows int the address table
func GetTopMinerChart() ([]*DBMinerRankInfo, error) {
	var topMinerInfo []*DBMinerRankInfo
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).All(&topMinerInfo)
	}
	err := withCollection(chartTopMinerRankTbl, query)
	return topMinerInfo, err
}

//AddNodeInfo add node info into database
func AddNodeInfo(nodeInfo *DBNodeInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(nodeInfo)
	}
	err := withCollection(nodeInfoTbl, query)
	return err
}

//DeleteNodeInfo delete node info from database
func DeleteNodeInfo(nodeInfo *DBNodeInfo) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"host": nodeInfo.Host})
	}
	err := withCollection(nodeInfoTbl, query)
	return err
}

//GetNodeInfo get node info from database
func GetNodeInfo(host string) (*DBNodeInfo, error) {
	dbNodeInfo := new(DBNodeInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"host": host}).One(dbNodeInfo)
	}
	err := withCollection(nodeInfoTbl, query)
	return dbNodeInfo, err
}

//GetNodeInfoByID get node info from database by node id
func GetNodeInfoByID(id string) (*DBNodeInfo, error) {
	dbNodeInfo := new(DBNodeInfo)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"id": id}).One(dbNodeInfo)
	}
	err := withCollection(nodeInfoTbl, query)
	return dbNodeInfo, err
}

//GetNodeInfos get all node infos from database
func GetNodeInfos() ([]*DBNodeInfo, error) {
	var nodeInfos []*DBNodeInfo
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).All(&nodeInfos)
	}
	err := withCollection(nodeInfoTbl, query)
	return nodeInfos, err
}

//GetNodeCnt get row count of the node table
func GetNodeCnt() (uint64, error) {
	var NodeCnt uint64
	query := func(c *mgo.Collection) error {
		var err error
		//TODO: fix this overflow
		var temp int
		temp, err = c.Count()
		NodeCnt = uint64(temp)
		return err
	}
	err := withCollection(nodeInfoTbl, query)
	return NodeCnt, err
}
