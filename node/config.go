/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package node

import "time"

//Config server config
type Config struct {
	RPCNodes        []string
	WriteLog        bool
	LogLevel        string
	LogFile         string
	DataBaseConnURL string
	DataBaseName    string
	Interval        time.Duration
	ExpireTime      int64
}
