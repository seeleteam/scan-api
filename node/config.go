/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package node

import (
	"time"

	"github.com/seeleteam/scan-api/common"
)

//Config server config
type Config struct {
	RPCNodes   []string
	WriteLog   bool
	LogLevel   string
	LogFile    string
	DataBase   *common.DataBaseConfig
	Interval   time.Duration
	ExpireTime int64
}
