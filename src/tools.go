package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

const GET_VETH_INFO_SCRIPT string = "tools/get-veth-info.sh"

type veth_info struct {
	Name  string `json:"veth-name"`
	Index string `json:"veth-id"`
	MAC   string `json:"veth-mac"`
}

func extractVethIDFromContainerID(containerid string) (*veth_info, error) {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, GET_VETH_INFO_SCRIPT)

	// /bin/bash tools/get-veth-info.sh <container-id>
	stdout, err := exec.Command("/bin/bash", path, string(containerid)).Output()
	if err != nil {
		return nil, err
	}

	info := new(veth_info)
	err = json.Unmarshal([]byte(stdout), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
