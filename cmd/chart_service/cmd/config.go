/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"scan-api/chart"
	"scan-api/database"
	"scan-api/rpc"
)

var (
	configFile *string
)

// LoadConfigFromFile unmarshal config from a file
func LoadConfigFromFile(filepath string) (statistics.Config, error) {
	var config statistics.Config
	buff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(buff, &config)
	//set database
	if len(config.DataBaseConnURL) > 0 {
		database.ConnURL = config.DataBaseConnURL
	}

	if len(config.DataBaseName) > 0 {
		database.DataBaseName = config.DataBaseName
	}

	//set rpc
	if len(config.RPCURL) > 0 {
		rpc.RPCURL = config.RPCURL
	}
	return config, err
}
