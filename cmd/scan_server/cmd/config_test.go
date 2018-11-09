/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */
package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile_SingleMode(t *testing.T) {
	fp := `./testfile/server1_test.json`
	cfg, err := LoadConfigFromFile(fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.GinMode, "debug")
	assert.Equal(t, cfg.Addr, ":8888")
	assert.Equal(t, cfg.LimitConnection, 0)
	assert.Equal(t, cfg.DefaultHammerTime, time.Duration(30))
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "scan-api")
	assert.Equal(t, cfg.DisableConsoleColor, true)
	assert.Equal(t, cfg.LimitConnections, 0)
	assert.Equal(t, cfg.MaxHeaderBytes, uint(20))
	assert.Equal(t, cfg.ReadTimeout, time.Duration(300))
	assert.Equal(t, cfg.IdleTimeout, time.Duration(0))
	assert.Equal(t, cfg.WriteTimeout, time.Duration(120))
	assert.Equal(t, cfg.SyncSwitch, true)
	assert.Equal(t, cfg.BlockCacheLimit, 1024)
	assert.Equal(t, cfg.TransCacheLimit, 1024)
	assert.Equal(t, cfg.Interval, time.Duration(30))
	assert.Equal(t, cfg.DataBase.DataBaseMode, "single")
	assert.Equal(t, cfg.DataBase.DataBaseConnURLs, []string{"127.0.0.1:27017"})
	assert.Equal(t, cfg.DataBase.DataBaseName, "seele")
	assert.Equal(t, cfg.DataBase.UseAuthentication, false)
	assert.Equal(t, cfg.DataBase.User, "scan")
	assert.Equal(t, cfg.DataBase.Pwd, "123456")
}

func TestLoadConfigFromFile_ReplsetMode(t *testing.T) {
	fp := `./testfile/server2_test.json`
	cfg, err := LoadConfigFromFile(fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.GinMode, "debug")
	assert.Equal(t, cfg.Addr, ":8888")
	assert.Equal(t, cfg.LimitConnection, 0)
	assert.Equal(t, cfg.DefaultHammerTime, time.Duration(30))
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "scan-api")
	assert.Equal(t, cfg.DisableConsoleColor, true)
	assert.Equal(t, cfg.LimitConnections, 0)
	assert.Equal(t, cfg.MaxHeaderBytes, uint(20))
	assert.Equal(t, cfg.ReadTimeout, time.Duration(300))
	assert.Equal(t, cfg.IdleTimeout, time.Duration(0))
	assert.Equal(t, cfg.WriteTimeout, time.Duration(120))
	assert.Equal(t, cfg.SyncSwitch, true)
	assert.Equal(t, cfg.BlockCacheLimit, 1024)
	assert.Equal(t, cfg.TransCacheLimit, 1024)
	assert.Equal(t, cfg.Interval, time.Duration(30))
	assert.Equal(t, cfg.DataBase.DataBaseMode, "replset")
	assert.Equal(t, cfg.DataBase.DataBaseReplsetName, "scan")
	assert.Equal(t, cfg.DataBase.DataBaseConnURLs, []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"})
	assert.Equal(t, cfg.DataBase.DataBaseName, "seele")
	assert.Equal(t, cfg.DataBase.UseAuthentication, false)
	assert.Equal(t, cfg.DataBase.User, "scan")
	assert.Equal(t, cfg.DataBase.Pwd, "123456")
}
