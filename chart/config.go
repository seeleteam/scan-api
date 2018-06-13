/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package chart

import "sync"

//Config server config
type Config struct {
	RPCURL          string
	WriteLog        bool
	LogLevel        string
	LogFile         string
	DataBaseConnURL string
	DataBaseName    string
	ShardCount      int
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
