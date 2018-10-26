/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package syncer

import (
	"sort"
	"testing"

	"github.com/seeleteam/scan-api/database"
	"github.com/stretchr/testify/assert"
)

type dates []string

func (c dates) Len() int           { return len(c) }
func (c dates) Less(i, j int) bool { return c[i] < c[j] }
func (c dates) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func TestGetTxDate(t *testing.T) {
	txs := []*database.DBTx{
		&database.DBTx{
			Timetxs: "2018-10-11",
		},
		&database.DBTx{
			Timetxs: "2018-10-12",
		},
		&database.DBTx{
			Timetxs: "2018-10-13",
		},
		&database.DBTx{
			Timetxs: "2018-10-14",
		},
		&database.DBTx{
			Timetxs: "2018-10-12",
		},
		&database.DBTx{
			Timetxs: "2018-10-12",
		},
	}
	date := dates(getTxDate(txs))
	sort.Sort(date)
	assert.Equal(t, []string(date), []string{"2018-10-11", "2018-10-12", "2018-10-13", "2018-10-14"})
}

func TestFilterDate(t *testing.T) {
	dates := []string{"2018-10-11", "2018-10-12", "2018-10-13", "2018-10-14"}
	limit := "2018-10-12"
	dates = filterDate(dates, limit)
	assert.Equal(t, dates, []string{"2018-10-12", "2018-10-13", "2018-10-14"})
}

func TestNextDate(t *testing.T) {
	date := "2018-10-11"
	nextDate(&date)
	assert.Equal(t, date, "2018-10-12")
}
