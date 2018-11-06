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
	assert.Equal(t, cfg.RpcURL, "127.0.0.1:55027")
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "seele-syncer")
	assert.Equal(t, cfg.DataBaseMode, "single")
	assert.Equal(t, cfg.DataBaseConnURL, []string{"127.0.0.1:27017"})
	assert.Equal(t, cfg.DataBaseName, "seele")
	assert.Equal(t, cfg.SyncInterval, time.Duration(3))
	assert.Equal(t, cfg.ShardNumber, 1)
}

func TestLoadConfigFromFile_ReplsetMode(t *testing.T) {
	fp := `./testfile/server2_test.json`
	cfg, err := LoadConfigFromFile(fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.RpcURL, "127.0.0.1:55027")
	assert.Equal(t, cfg.WriteLog, true)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFile, "seele-syncer")
	assert.Equal(t, cfg.DataBaseMode, "replset")
	assert.Equal(t, cfg.DataBaseConnURL, []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"})
	assert.Equal(t, cfg.DataBaseName, "seele")
	assert.Equal(t, cfg.SyncInterval, time.Duration(3))
	assert.Equal(t, cfg.ShardNumber, 1)
}
