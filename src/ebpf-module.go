package main

import (
	//	"io"
	//	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	v1 "k8s.io/api/core/v1"
	klog "k8s.io/klog/v2"
)

type ebpfModule struct {
	toolsPath string
	config    *EbpfPrograms
}

func (m *ebpfModule) getRequestedArg(requested string, info *veth_info) (string, error) {
	var arg string
	switch requested {
	case BPF_ARG_VETH_NAME:
		arg = info.VethName
	case BPF_ARG_VETH_ID:
		arg = info.VethIndex
	case BPF_ARG_VETH_MAC:
		arg = info.VethMac
	case BPF_ARG_VPEER_NAME:
		arg = info.VpeerName
	case BPF_ARG_VPEER_ID:
		arg = info.VpeerIndex
	case BPF_ARG_VPEER_MAC:
		arg = info.VpeerMac
	case BPF_ARG_NAMESPACE:
		arg = info.Namespace
	default:
		return "", fmt.Errorf("Requested arg %s not available in veth-info.", requested)
	}
	return arg, nil
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

	klog.Infof("veth info extracted")

	// Check the requested attachment
	cmd := exec.Command("/bin/bash", bpfProgram.CMD)
	/*
		stderr, err := cmd.StderrPipe()
		if err != nil {
			klog.Errorf("could not get stderr pipe: %v", err)
		} else {
			go func() {
				var cmdStr string = cmd.String()
				reader := io.Reader(stderr)
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					msg := scanner.Text()
					fmt.Printf("stderr[%s]: %s\n",cmdStr, msg)
				}
			}()
		}
	*/
	cmd.Dir = bpfProgram.Path
	for _, requested := range bpfProgram.Env {
		param, err := m.getRequestedArg(requested, info)
		if err != nil {
			return err
		}
		envArg := requested + "=" + param
		cmd.Env = append(cmd.Env, envArg)
	}

	klog.Infof("dir is %s", bpfProgram.Path)
	klog.Infof("env is %v", cmd.Env)
	klog.Infof("cmd is %s", cmd)

	stdout, err := cmd.Output()
	klog.Infof("stdout is :-\n\t" + string(stdout))
	if err != nil {
		return err
	}
	klog.Infof("command executed")
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
		info, err = extractVethIDFromContainerID(container_id, m.toolsPath)
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
