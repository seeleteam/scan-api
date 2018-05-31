/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"strconv"
	"time"
)

const (
	twelveHours = 12 * 60 * 60
)

var (
	//Avg12HoursHashrate average hashrate in the past 12 hours
	Avg12HoursHashrate float64
)

//ProcessLast12HoursHashRate Calculate hashrate in the past 12 hours
func ProcessLast12HoursHashRate() {
	now := time.Now()
	last12HoursTime := now.Add(-time.Hour * 12)
	var dbBlocks []*DBBlock
	dbBlocks, err := GetBlocksByTime(last12HoursTime.Unix(), now.Unix())
	if err != nil {
		return
	}

	var difficulty uint64

	for i := 0; i < len(dbBlocks); i++ {
		diff, err := strconv.ParseUint(dbBlocks[i].Difficulty, 10, 64)
		if err == nil {
			difficulty += diff
		}
	}

	Avg12HoursHashrate = float64(difficulty) / twelveHours
	return
}
