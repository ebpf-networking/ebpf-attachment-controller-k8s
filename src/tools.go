package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
)

const GET_VETH_INFO_SCRIPT string = "get-veth-info.sh"

type veth_info struct {
	Name  string `json:"veth-name"`
	Index string `json:"veth-id"`
	MAC   string `json:"veth-mac"`
}

func extractVethIDFromContainerID(containerid string, toolsPath string) (*veth_info, error) {

	cmd := exec.Command("/bin/bash", GET_VETH_INFO_SCRIPT, string(containerid))
	cmd.Dir = toolsPath

	// /bin/bash tools/get-veth-info.sh <container-id>
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Cmd executed. output %s",string(stdout))

	info := new(veth_info)
	err = json.Unmarshal([]byte(stdout), &info)
	if err != nil {
		fmt.Printf("couldn't unmarshall json")
		return nil, err
	}
	return info, nil
}
