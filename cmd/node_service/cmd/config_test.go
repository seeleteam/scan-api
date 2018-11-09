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
	assert.Equal(t, cfg.RPCNodes, []string{"106.75.80.93:55026"})
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "node_service")
	assert.Equal(t, cfg.DataBase.DataBaseMode, "single")
	assert.Equal(t, cfg.DataBase.DataBaseConnURLs, []string{"127.0.0.1:27017"})
	assert.Equal(t, cfg.DataBase.DataBaseName, "seele")
	assert.Equal(t, cfg.DataBase.UseAuthentication, false)
	assert.Equal(t, cfg.DataBase.User, "scan")
	assert.Equal(t, cfg.DataBase.Pwd, "123456")
	assert.Equal(t, cfg.Interval, time.Duration(60))
	assert.Equal(t, cfg.ExpireTime, int64(60))
}

func TestLoadConfigFromFile_ReplsetMode(t *testing.T) {
	fp := `./testfile/server2_test.json`
	cfg, err := LoadConfigFromFile(fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.RPCNodes, []string{"106.75.80.93:55026"})
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "node_service")
	assert.Equal(t, cfg.DataBase.DataBaseMode, "replset")
	assert.Equal(t, cfg.DataBase.DataBaseReplsetName, "scan")
	assert.Equal(t, cfg.DataBase.DataBaseConnURLs, []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"})
	assert.Equal(t, cfg.DataBase.DataBaseName, "seele")
	assert.Equal(t, cfg.DataBase.UseAuthentication, false)
	assert.Equal(t, cfg.DataBase.User, "scan")
	assert.Equal(t, cfg.DataBase.Pwd, "123456")
	assert.Equal(t, cfg.Interval, time.Duration(60))
	assert.Equal(t, cfg.ExpireTime, int64(60))
}
