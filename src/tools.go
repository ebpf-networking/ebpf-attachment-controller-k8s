package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	klog "k8s.io/klog/v2"
)

const GET_VETH_INFO_SCRIPT string = "get-veth-info.sh"

type veth_info struct {
	VethName   string `json:"veth-name"`
	VethIndex  string `json:"veth-id"`
	VethMac    string `json:"veth-mac"`
	VpeerName  string `json:"vpeer-name"`
	VpeerIndex string `json:"vpeer-id"`
	VpeerMac   string `json:"vpeer-mac"`
	Namespace  string `json:"netns"`
}

func extractVethIDFromContainerID(containerid string, toolsPath string) (*veth_info, error) {

	cmd := exec.Command("/bin/bash", GET_VETH_INFO_SCRIPT, string(containerid))
	cmd.Dir = toolsPath

	// /bin/bash tools/get-veth-info.sh <container-id>
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	klog.Infof("Tool %s executed. output:-\n\n%s\n", cmd, string(stdout))

	info := new(veth_info)
	err = json.Unmarshal([]byte(stdout), &info)
	if err != nil {
		fmt.Printf("couldn't unmarshall json")
		return nil, err
	}
	return info, nil
}
