package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func RunConfig() {
	data, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/event/user/testing/config/config.json")
	if err != nil {
		panic(err)
	}
	var variables map[string]string
	if err := json.Unmarshal(data, &variables); err != nil {
		panic(err)
	}
	for key, value := range variables {
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				panic(err)
			}
		}

	}
}
