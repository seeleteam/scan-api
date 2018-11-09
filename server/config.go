/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package server

import (
	"time"

	"github.com/seeleteam/scan-api/common"
)

//Config server config
type Config struct {
	GinMode             string
	Addr                string
	LimitConnection     int
	DefaultHammerTime   time.Duration
	WriteLog            bool
	LogLevel            string
	LogFile             string
	DisableConsoleColor bool
	LimitConnections    int
	MaxHeaderBytes      uint
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	SyncSwitch          bool
	BlockCacheLimit     int
	TransCacheLimit     int
	DataBase            *common.DataBaseConfig
	Interval            time.Duration
}
