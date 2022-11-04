package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	VETH_NAME string = "VETH_NAME"
	VETH_ID          = "VETH_ID"
	VETH_MAC         = "VETH_MAC"
)

type BPFProgram struct {
	Prog string   `json:"program"` // Program name which is requested by
	CMD  string   `json:"cmd"`     // bash script to execute fror the program
	Path string   `json:"path"`    // Path (including script name) to the script which needs to be executed
	Env  []string `json:"args"`    // The arguments needed by the BPF program
}

type EbpfPrograms struct {
	Programs map[string]BPFProgram `json:"attachments"`
}

func loadConfigFromPath(configpath string) (*EbpfPrograms, error) {
	config := new(EbpfPrograms)
	err := loadJSONFromFile(configpath, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadJSONFromFile(filename string, v interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s" + err.Error())
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json" + err.Error())
	}
	return nil
}
