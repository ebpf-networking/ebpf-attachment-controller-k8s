package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	v1 "k8s.io/api/core/v1"
	klog "k8s.io/klog/v2"
)

type ebpfModule struct {
	config *EbpfPrograms
}

func (m *ebpfModule) AttachEBPFNetwork(pod *v1.Pod, program string) error {

	bpfProgram, ok := m.config.Programs[program]
	if !ok {
		klog.Infof("Requested EBPF program %s not found", program)
		return fmt.Errorf("Requested EBPF program %s not found", program)
	}

	podIP := pod.Status.PodIP
	hostIP := pod.Status.HostIP

	klog.Infof("pod name:%s, ip:%s, hostip:%s", pod.Name, podIP, hostIP)

	info, err := m.extractVethInfo(pod)
	if err != nil {
		return err
	}

	// Check the requested attachment
	cmd := exec.Command("/bin/bash", bpfProgram.CMD)
	cmd.Dir = bpfProgram.Path

	for _, requested := range bpfProgram.Env {
		switch requested {
		case VETH_NAME:
			cmd.Env = append(cmd.Env, VETH_NAME+"="+info.Name)
		case VETH_ID:
			cmd.Env = append(cmd.Env, VETH_ID+"="+info.Index)
		case VETH_MAC:
			cmd.Env = append(cmd.Env, VETH_MAC+"="+info.MAC)
		}
	}

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	klog.Infof(string(stdout))
	return nil
}

func (m *ebpfModule) DeleteEBPFNetwork(pod *v1.Pod, program string) {
	klog.Infof("deleteEBPFNetwork Not implemented")
}

func (m *ebpfModule) extractVethInfo(pod *v1.Pod) (*veth_info, error) {
	var info *veth_info = nil
	var err error

	// Extraction from any one container is fine.
	// All containers of a pod are on same host with same namespace.
	for i, containerstatus := range pod.Status.ContainerStatuses {
		container := containerstatus.ContainerID
		container_id := strings.TrimPrefix(container, "containerd://")
		klog.Infof("container[%d] : %s", i, container_id)

		// Get the veth information from container id
		info, err = extractVethIDFromContainerID(container_id)
		if err != nil {
			klog.Errorf(err.Error())
			continue
		}
		// extracted
		break
	}
	if info == nil {
		return nil, errors.New("failed to get veth info from any container")
	}
	return info, nil
}
