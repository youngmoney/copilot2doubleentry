package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TransactionMatch struct {
	// TODO: add regex to the json names
	DescriptionRegex string `json:"description"`
	AccountRgex      string `json:"account"`
	// if positive in copilot
	Outgoing bool `json:"outgoing"`
}

type Override struct {
	Match        TransactionMatch `json:"match"`
	Account      string           `json:"account"`
	SplitAccount string           `json:"split_account"`
	AlwaysPair   bool             `json:"always_pair"`
}

type Overrides struct {
	Income   []Override `json:"income"`
	Expense  []Override `json:"expense"`
	Transfer []Override `json:"transfer"`
}

type Config struct {
	Overrides Overrides `json:"overrides"`
}

func ReadConfig(filename string) Config {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := Config{}
	if err := json.Unmarshal(raw, &config); err != nil {
		panic(err)
	}

	fmt.Println(config)

	return config
}
