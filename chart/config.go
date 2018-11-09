/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package chart

import (
	"sync"

	"github.com/seeleteam/scan-api/common"
)

//Config server config
type Config struct {
	RPCURL     string
	WriteLog   bool
	LogLevel   string
	LogFile    string
	DataBase   *common.DataBaseConfig
	ShardCount int
}

//ProcessFunc ChartProcessFunc is entrance of the chart service needed to be start
type ProcessFunc func(wg *sync.WaitGroup)

var (
	//ProcessFuncs chart processors
	processFuncs []ProcessFunc
)

//RegisterProcessFunc register an process func into chart service
func RegisterProcessFunc(processFunc ProcessFunc) {
	processFuncs = append(processFuncs, processFunc)
}

//GetProcessFuncs return process func slice
func GetProcessFuncs() []ProcessFunc {
	return processFuncs
}
