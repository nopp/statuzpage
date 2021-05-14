package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	StatuzpageAPI string `json:"statuzpage-api"`
	MySQLHost     string `json:"mysql-host"`
	MySQLUser     string `json:"mysql-user"`
	MySQLPass     string `json:"mysql-password"`
	MySQLDb       string `json:"mysql-db"`
	Token         string `json:"token"`
}

// LoadConfiguration from config.json
func LoadConfiguration() config {

	var config config

	configFile, err := ioutil.ReadFile("/etc/statuzpage-agent/config.json")
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(configFile, &config)
	return config
}
