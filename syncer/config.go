package syncer

import (
	"time"
)

// DataBaseConfig database info
type DataBaseConfig struct {
	DataBaseMode    string
	DataBaseReplsetName string
	DataBaseConnURLs []string
	DataBaseName    string
	UseAuthentication bool
	User string
	Pwd string
}

func (db *DataBaseConfig) GetDBName() string {return db.DataBaseName}
func (db *DataBaseConfig) GetDBMode() string {return db.DataBaseMode}
func (db *DataBaseConfig) GetReplsetName() string {return db.DataBaseReplsetName}
func (db *DataBaseConfig) GetConnURLs() []string {return db.DataBaseConnURLs}
func (db *DataBaseConfig) GetUseAuthentication() bool {return db.UseAuthentication}
func (db *DataBaseConfig) GetUser() string {return db.User}
func (db *DataBaseConfig) GetPwd() string {return db.Pwd}

// Config server config
type Config struct {
	RpcURL       string
	WriteLog     bool
	LogLevel     string
	LogFile      string
	DataBase     *DataBaseConfig
	SyncInterval time.Duration
	ShardNumber  int
}
