package main

import (
	"encoding/json"
	"io/ioutil"
)

func main() {
	cfg, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}
	var rawcfg RawServerConfig
	err = json.Unmarshal(cfg, &rawcfg)
	if err != nil {
		panic(err.Error())
	}
	var serverConfig = Init(&rawcfg)
	HandleTCP(serverConfig)
}
