/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/seeleteam/scan-api/server"
)

// LoadConfigFromFile unmarshal config from a file
func LoadConfigFromFile(filepath string) (server.Config, error) {
	var config server.Config
	buff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(buff, &config)

	return config, err
}
