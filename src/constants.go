package main

const (
	BPF_ARG_VETH_NAME string = "VETH_NAME"
	BPF_ARG_VETH_ID          = "VETH_ID"
	BPF_ARG_VETH_MAC         = "VETH_MAC"
)

const (
	POD_ANNOTATION_EBPF_ATTACHMENT   string = "ebpf-attachment"
	POD_ANNOTATION_ATTACHMENT_STATUS string = "ebpf-attachment-status"
)

const (
	EBPF_ATTACHED string = "attached"
	EBPF_FAILED   string = "failed"
)

const (
	ENV_MY_HOST               string = "MY_NODE_NAME"
	ENV_KUBECONFIG            string = "KUBECONFIG"
	ENV_CONFIG_PATH           string = "CONTROLLER_CONFIG"
	ENV_CONTROLLER_TOOLS_PATH string = "CONTROLLER_TOOLS"
)

const (
	DEFAULT_CONFIG_PATH string = "/opt/config/controller-config.json"
	DEFAULT_TOOLS_PATH  string = "/opt/src/tools/"
)
