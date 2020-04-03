package shm

import (
	"encoding/json"
	"io/ioutil"
)

var config map[string]interface{}

func ConfigGet(key string) (interface{}, error) {

	if config != nil {
		return config[key], nil
	}
	content, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return config[key], nil
}
