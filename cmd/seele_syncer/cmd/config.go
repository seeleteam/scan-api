package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/seeleteam/scan-api/syncer"
)

// LoadConfigFromFile unmarshal config from a file
func LoadConfigFromFile(filepath string) (syncer.Config, error) {
	var config syncer.Config
	buff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(buff, &config)
	return config, err
}
