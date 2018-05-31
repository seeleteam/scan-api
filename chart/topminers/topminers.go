/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package topminers

import (
	"scan-api/database"
	"sort"
	"sync"
	"time"
)

const (
	sevenDays = time.Hour * 24 * 7
)

//MinerRankInfoSlice rank array
type MinerRankInfoSlice []database.DBSingleMinerRankInfo

func (s MinerRankInfoSlice) Len() int           { return len(s) }
func (s MinerRankInfoSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MinerRankInfoSlice) Less(i, j int) bool { return s[i].Mined > s[j].Mined }

//Process set an timer to calculate top miners every day
func Process(wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessTopMiners()

	for {
		now := time.Now()
		// calcuate next zero hour
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//calcuate last day transactions
		ProcessTopMiners()
	}
}

//ProcessTopMiners get the top miners who mined the most blocks in the last 7 days
func ProcessTopMiners() bool {
	now := time.Now()

	thisZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastDayZeroTime := thisZeroTime.Add(-sevenDays)
	var dbBlocks []*database.DBBlock
	var err error
	dbBlocks, err = database.GetBlocksByTime(lastDayZeroTime.Unix(), thisZeroTime.Unix())
	if err != nil {
		return false
	}

	miners := make(map[string]int)
	for i := 0; i < len(dbBlocks); i++ {
		block := dbBlocks[i]
		miners[block.Creator]++
	}

	var TopSevenDaysMiners MinerRankInfoSlice
	for k, v := range miners {
		miner := database.DBSingleMinerRankInfo{
			Address:    k,
			Mined:      v,
			Percentage: float64(v) / float64(len(dbBlocks)),
		}
		TopSevenDaysMiners = append(TopSevenDaysMiners, miner)
	}

	sort.Stable(TopSevenDaysMiners)
	dbRank := database.DBMinerRankInfo{Rank: TopSevenDaysMiners}
	database.RemoveTopMinerInfo()
	database.AddTopMinerInfo(&dbRank)

	return true
}
