/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package server

import (
	"time"
)

//Config server config
type Config struct {
	GinMode             string
	Addr                string
	LimitConnection     int
	DefaultHammerTime   time.Duration
	RPCURL              string
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
	DataBaseConnURL     string
	DataBaseName        string
	Interval            time.Duration
}
