package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Root string `json:"root"`
	Walk string `json:"walk"`
}

func readConfig() *config {
	bytes, err := ioutil.ReadFile("kfs.json")
	if err != nil {
		panic(err)
	}
	config := &config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		panic(err)
	}
	return config
}
