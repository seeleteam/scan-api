/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seeleteam/scan-api/chart/topminers"
	"github.com/seeleteam/scan-api/database"
)

//Set an simple set to help calc chart infomation
type Set struct {
	m map[int64]bool
}

const showLimit = 15

//NewSet return an new set
func NewSet() *Set {
	return &Set{
		m: map[int64]bool{},
	}
}

//Add add an item into set
func (s *Set) Add(item int64) {
	s.m[item] = true
}

//Remove remove an item from set
func (s *Set) Remove(item int64) {
	delete(s.m, item)
}

//Has check the item is exist in the set
func (s *Set) Has(item int64) bool {
	_, ok := s.m[item]
	return ok
}

//ChartHandler handle all chart request
type ChartHandler struct {
	DBClient ChartInfoDB
}

//addUpOneDayTrans
func addUpOneDayTrans(oneDayTrans []*database.DBOneDayTxInfo) []*database.DBOneDayTxInfo {
	var ret []*database.DBOneDayTxInfo
	set := NewSet()
	for i := 0; i < len(oneDayTrans); i++ {
		if set.Has(oneDayTrans[i].TimeStamp) {
			continue
		}

		var info database.DBOneDayTxInfo
		set.Add(oneDayTrans[i].TimeStamp)
		for j := 0; j < len(oneDayTrans); j++ {
			if oneDayTrans[j].TimeStamp == oneDayTrans[i].TimeStamp {
				info.TimeStamp = oneDayTrans[j].TimeStamp
				info.ShardNumber = 1
				info.TotalBlocks += oneDayTrans[j].TotalBlocks
				info.TotalTxs += oneDayTrans[j].TotalTxs
			}
		}

		ret = append(ret, &info)
	}

	return ret
}

//addUpOneDayAddressInfo
func addUpOneDayAddressInfo(oneDayAddressInfos []*database.DBOneDayAddressInfo) []*database.DBOneDayAddressInfo {
	var ret []*database.DBOneDayAddressInfo
	set := NewSet()
	for i := 0; i < len(oneDayAddressInfos); i++ {
		if set.Has(oneDayAddressInfos[i].TimeStamp) {
			continue
		}

		var info database.DBOneDayAddressInfo
		set.Add(oneDayAddressInfos[i].TimeStamp)
		for j := 0; j < len(oneDayAddressInfos); j++ {
			if oneDayAddressInfos[j].TimeStamp == oneDayAddressInfos[i].TimeStamp {
				info.TimeStamp = oneDayAddressInfos[j].TimeStamp
				info.ShardNumber = 1
				info.TodayIncrease += oneDayAddressInfos[j].TodayIncrease
				info.TotalAddresss += oneDayAddressInfos[j].TotalAddresss
			}
		}

		ret = append(ret, &info)
	}
	return ret
}

//addUpOneDayBlockInfo
func addUpOneDayBlockInfo(oneDayBlocks []*database.DBOneDayBlockInfo) []*database.DBOneDayBlockInfo {
	var ret []*database.DBOneDayBlockInfo
	set := NewSet()
	for i := 0; i < len(oneDayBlocks); i++ {
		if set.Has(oneDayBlocks[i].TimeStamp) {
			continue
		}

		var info database.DBOneDayBlockInfo
		set.Add(oneDayBlocks[i].TimeStamp)
		for j := 0; j < len(oneDayBlocks); j++ {
			if oneDayBlocks[j].TimeStamp == oneDayBlocks[i].TimeStamp {
				info.TimeStamp = oneDayBlocks[j].TimeStamp
				info.ShardNumber = 1
				info.TotalBlocks += oneDayBlocks[j].TotalBlocks
				info.Rewards += oneDayBlocks[j].Rewards
			}
		}

		ret = append(ret, &info)
	}
	return ret
}

func averageOneDayBlockDifficulty(oneDayBlockDifficulties []*database.DBOneDayBlockDifficulty) []*database.DBOneDayBlockDifficulty {
	var ret []*database.DBOneDayBlockDifficulty
	set := NewSet()
	for i := 0; i < len(oneDayBlockDifficulties); i++ {
		if set.Has(oneDayBlockDifficulties[i].TimeStamp) {
			continue
		}

		countShard := 0
		var info database.DBOneDayBlockDifficulty
		set.Add(oneDayBlockDifficulties[i].TimeStamp)
		for j := 0; j < len(oneDayBlockDifficulties); j++ {
			if oneDayBlockDifficulties[j].TimeStamp == oneDayBlockDifficulties[i].TimeStamp {
				info.TimeStamp = oneDayBlockDifficulties[j].TimeStamp
				info.ShardNumber = 1
				info.Difficulty += oneDayBlockDifficulties[j].Difficulty
				countShard++
			}
		}

		info.Difficulty = info.Difficulty / float64(countShard)
		ret = append(ret, &info)
	}
	return ret
}

func averageOneDayBlockTime(oneDayBlockTimes []*database.DBOneDayBlockAvgTime) []*database.DBOneDayBlockAvgTime {
	var ret []*database.DBOneDayBlockAvgTime
	set := NewSet()
	for i := 0; i < len(oneDayBlockTimes); i++ {
		if set.Has(oneDayBlockTimes[i].TimeStamp) {
			continue
		}

		countShard := 0
		var info database.DBOneDayBlockAvgTime
		set.Add(oneDayBlockTimes[i].TimeStamp)
		for j := 0; j < len(oneDayBlockTimes); j++ {
			if oneDayBlockTimes[j].TimeStamp == oneDayBlockTimes[i].TimeStamp {
				info.TimeStamp = oneDayBlockTimes[j].TimeStamp
				info.ShardNumber = 1
				info.AvgTime += oneDayBlockTimes[j].AvgTime
				countShard++
			}
		}

		info.AvgTime = info.AvgTime / float64(countShard)
		ret = append(ret, &info)
	}
	return ret
}

func averageOneDayHashRate(oneDayHashRates []*database.DBOneDayHashRate) []*database.DBOneDayHashRate {
	var ret []*database.DBOneDayHashRate
	set := NewSet()
	for i := 0; i < len(oneDayHashRates); i++ {
		if set.Has(oneDayHashRates[i].TimeStamp) {
			continue
		}

		countShard := 0
		var info database.DBOneDayHashRate
		set.Add(oneDayHashRates[i].TimeStamp)
		for j := 0; j < len(oneDayHashRates); j++ {
			if oneDayHashRates[j].TimeStamp == oneDayHashRates[i].TimeStamp {
				info.TimeStamp = oneDayHashRates[j].TimeStamp
				info.ShardNumber = 1
				info.HashRate += oneDayHashRates[j].HashRate
				countShard++
			}
		}

		info.HashRate = info.HashRate / float64(countShard)
		ret = append(ret, &info)
	}
	return ret
}

//GetTxHistory handler for transaction history chart
func (h *ChartHandler) GetTxHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbClinet := h.DBClient

		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if s < 0 {
			s = 1
		}

		shardNumber := int(s)
		if shardNumber == 0 {
			//get all info
			oneDayTrans, err := dbClinet.GetTransInfoChart()
			if err != nil {
				responseError(c, errGetTxChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayTrans = addUpOneDayTrans(oneDayTrans)
			oneDayAddresses, err := dbClinet.GetOneDayAddressesChart()
			if err != nil {
				responseError(c, errGetAddressChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayAddresses = addUpOneDayAddressInfo(oneDayAddresses)

			oneDayBlocks, err := dbClinet.GetOneDayBlocksChart()
			if err != nil {
				responseError(c, errGetBlockCountAndRewardChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayBlocks = addUpOneDayBlockInfo(oneDayBlocks)

			//all the tables should have the same timeline
			if len(oneDayTrans) != len(oneDayAddresses) ||
				len(oneDayAddresses) != len(oneDayBlocks) {
				responseError(c, errDBDataError, http.StatusInternalServerError, apiInternalError)
				return
			}

			var retTxHistorys []RetOneDayTxInfo
			for i := 0; i < len(oneDayTrans); i++ {
				retTx := RetOneDayTxInfo{
					TotalTxs:      oneDayTrans[i].TotalTxs,
					TotalBlocks:   int(oneDayBlocks[i].TotalBlocks),
					Rewards:       oneDayBlocks[i].Rewards,
					TotalAddresss: oneDayAddresses[i].TotalAddresss,
					TodayIncrease: oneDayAddresses[i].TodayIncrease,
					TimeStamp:     oneDayTrans[i].TimeStamp,
				}
				retTxHistorys = append(retTxHistorys, retTx)
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    retTxHistorys,
			})
		} else {
			//get shader info
			oneDayTrans, err := dbClinet.GetTransInfoChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetTxChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayAddresses, err := dbClinet.GetOneDayAddressesChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetAddressChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayBlockDifficulties, err := dbClinet.GetOneDayBlockDifficultyChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockDifficultyChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayBlockTimes, err := dbClinet.GetOneDayBlockAvgTimeChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockTimeChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayBlocks, err := dbClinet.GetOneDayBlocksChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockCountAndRewardChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			oneDayHashRates, err := dbClinet.GetHashRateChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetHashRateChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			//all the tables should have the same timeline
			if len(oneDayTrans) != len(oneDayAddresses) ||
				len(oneDayAddresses) != len(oneDayBlockDifficulties) ||
				len(oneDayBlockDifficulties) != len(oneDayBlockTimes) ||
				len(oneDayBlockTimes) != len(oneDayBlocks) ||
				len(oneDayBlocks) != len(oneDayHashRates) {
				responseError(c, errDBDataError, http.StatusInternalServerError, apiInternalError)
				return
			}

			var retTxHistorys []RetOneDayTxInfo
			for i := 0; i < len(oneDayTrans); i++ {
				retTx := RetOneDayTxInfo{
					TotalTxs:      oneDayTrans[i].TotalTxs,
					TotalBlocks:   int(oneDayBlocks[i].TotalBlocks),
					HashRate:      oneDayHashRates[i].HashRate,
					Difficulty:    oneDayBlockDifficulties[i].Difficulty,
					AvgTime:       oneDayBlockTimes[i].AvgTime,
					Rewards:       oneDayBlocks[i].Rewards,
					TotalAddresss: oneDayAddresses[i].TotalAddresss,
					TodayIncrease: oneDayAddresses[i].TodayIncrease,
					TimeStamp:     oneDayTrans[i].TimeStamp,
				}
				retTxHistorys = append(retTxHistorys, retTx)
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    retTxHistorys,
			})
		}
	}
}

//GetEveryDayAddress handler for address chart
func (h *ChartHandler) GetEveryDayAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if s < 0 {
			s = 1
		}
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			oneDayAddresses, err := dbClinet.GetOneDayAddressesChart()
			if err != nil {
				responseError(c, errGetAddressChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayAddresses = addUpOneDayAddressInfo(oneDayAddresses)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayAddresses,
			})
		} else {
			oneDayAddresses, err := dbClinet.GetOneDayAddressesChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetAddressChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayAddresses,
			})
		}
	}
}

//GetEveryDayBlockDifficulty handler for block difficulty chart
func (h *ChartHandler) GetEveryDayBlockDifficulty() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			oneDayBlockDifficulties, err := dbClinet.GetOneDayBlockDifficultyChart()
			if err != nil {
				responseError(c, errGetBlockDifficultyChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayBlockDifficulties = averageOneDayBlockDifficulty(oneDayBlockDifficulties)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlockDifficulties,
			})
		} else {
			oneDayBlockDifficulties, err := dbClinet.GetOneDayBlockDifficultyChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockDifficultyChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlockDifficulties,
			})
		}

	}
}

//GetEveryDayBlockTime handler for avg block time chart
func (h *ChartHandler) GetEveryDayBlockTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			oneDayBlockTimes, err := dbClinet.GetOneDayBlockAvgTimeChart()
			if err != nil {
				responseError(c, errGetBlockTimeChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayBlockTimes = averageOneDayBlockTime(oneDayBlockTimes)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlockTimes,
			})
		} else {
			oneDayBlockTimes, err := dbClinet.GetOneDayBlockAvgTimeChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockTimeChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlockTimes,
			})
		}

	}
}

//GetEveryDayBlock handler for every day block chart
func (h *ChartHandler) GetEveryDayBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		if s < 0 {
			s = 1
		}
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			//get all info
			oneDayBlocks, err := dbClinet.GetOneDayBlocksChart()
			if err != nil {
				responseError(c, errGetBlockCountAndRewardChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayBlocks = addUpOneDayBlockInfo(oneDayBlocks)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlocks,
			})
		} else {
			//get shard info
			oneDayBlocks, err := dbClinet.GetOneDayBlocksChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetBlockCountAndRewardChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayBlocks,
			})
		}
	}
}

//GetTopMiners handler for every day block chart
func (h *ChartHandler) GetTopMiners() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			topMiners, err := dbClinet.GetTopMinerChart()
			if err != nil {
				responseError(c, errGetTopMinerChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			allMined := 0
			var TopMiners topminers.MinerRankInfoSlice
			for i := 0; i < len(topMiners); i++ {
				for j := 0; j < len(topMiners[i].Rank); j++ {
					TopMiners = append(TopMiners, topMiners[i].Rank[j])
					allMined += topMiners[i].Rank[j].Mined
				}
			}

			for i := 0; i < len(TopMiners); i++ {
				TopMiners[i].Percentage = float64(float64(TopMiners[i].Mined) / float64(allMined))
			}

			sort.Stable(TopMiners)
			if len(TopMiners) > showLimit {
				TopMiners = TopMiners[:showLimit]
			}

			topMiners = topMiners[:0]
			topMiners = append(topMiners, &database.DBMinerRankInfo{Rank: TopMiners, ShardNumber: 1})

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    topMiners,
			})
		} else {
			topMiners, err := dbClinet.GetTopMinerChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetTopMinerChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			if len(topMiners) > 0 {
				if len(topMiners[0].Rank) > showLimit {
					topMiners[0].Rank = topMiners[0].Rank[:showLimit]
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    topMiners,
			})
		}

	}
}

//GetEveryHashRate handler for every day hash Rate
func (h *ChartHandler) GetEveryHashRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := strconv.ParseInt(c.Query("s"), 10, 64)
		shardNumber := int(s)

		dbClinet := h.DBClient

		if shardNumber == 0 {
			oneDayHashRates, err := dbClinet.GetHashRateChart()
			if err != nil {
				responseError(c, errGetHashRateChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}
			oneDayHashRates = averageOneDayHashRate(oneDayHashRates)

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayHashRates,
			})
		} else {
			oneDayHashRates, err := dbClinet.GetHashRateChartByShardNumber(shardNumber)
			if err != nil {
				responseError(c, errGetHashRateChartError, http.StatusInternalServerError, apiDBQueryError)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    apiOk,
				"message": "",
				"data":    oneDayHashRates,
			})
		}
	}
}
