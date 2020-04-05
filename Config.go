package shm

import (
	"encoding/json"
	"io/ioutil"
)

var config map[string]interface{}

func GetConfig() (map[string]interface{}, error) {

	if config != nil {
		return config, nil
	}
	content, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
