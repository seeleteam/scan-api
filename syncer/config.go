package syncer

import (
	"time"
)

//Config server config
type Config struct {
	RpcURL          string
	WriteLog        bool
	LogLevel        string
	LogFile         string
	DataBaseConnURL string
	DataBaseName    string
	SyncInterval    time.Duration
	ShardNumber     int
}
