/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"net/http"
	"scan-api/database"

	"github.com/gin-gonic/gin"
)

//GetTxHistory handler for transaction history chart
func GetTxHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayTrans, err := database.GetTransInfoChart()
		if err != nil {
			responseError(c, errGetTxChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		oneDayAddresses, err := database.GetOneDayAddressesChart()
		if err != nil {
			responseError(c, errGetAddressChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		oneDayBlockDifficulties, err := database.GetOneDayBlockDifficultyChart()
		if err != nil {
			responseError(c, errGetBlockDifficultyChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		oneDayBlockTimes, err := database.GetOneDayBlockAvgTimeChart()
		if err != nil {
			responseError(c, errGetBlockTimeChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		oneDayBlocks, err := database.GetOneDayBlocksChart()
		if err != nil {
			responseError(c, errGetBlockCountAndRewardChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		oneDayHashRates, err := database.GetHashRateChart()
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

//GetEveryDayAddress handler for address chart
func GetEveryDayAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayAddresses, err := database.GetOneDayAddressesChart()
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

//GetEveryDayBlockDifficulty handler for block difficulty chart
func GetEveryDayBlockDifficulty() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayBlockDifficulties, err := database.GetOneDayBlockDifficultyChart()
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

//GetEveryDayBlockTime handler for avg block time chart
func GetEveryDayBlockTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayBlockTimes, err := database.GetOneDayBlockAvgTimeChart()
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

//GetEveryDayBlock handler for every day block chart
func GetEveryDayBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayBlocks, err := database.GetOneDayBlocksChart()
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

//GetTopMiners handler for every day block chart
func GetTopMiners() gin.HandlerFunc {
	return func(c *gin.Context) {
		topMiners, err := database.GetTopMinerChart()
		if err != nil {
			responseError(c, errGetTopMinerChartError, http.StatusInternalServerError, apiDBQueryError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    apiOk,
			"message": "",
			"data":    topMiners,
		})
	}
}

//GetEveryHashRate handler for every day hash Rate
func GetEveryHashRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		oneDayHashRates, err := database.GetHashRateChart()
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
