package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	klog "k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/util/logs"
)

const (
	MY_HOST             string = "MY_NODE_NAME"
	DEFAULT_CONFIG_PATH string = "/opt/config/controller-config.json"
)

type KubeClient struct {
	Client kubernetes.Interface
}

func (k8 *KubeClient) InitKubeClient() error {
	var kubeConfig *rest.Config
	var err error

	kubeConfigPath := os.Getenv("KUBECONFIG")
	if len(kubeConfigPath) > 0 {
		klog.Infof("loading from kubeconfig file %s", kubeConfigPath)
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return err
		}
	} else {
		klog.Infof("kubeconfig is null. loading incluster config")
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}
	klog.Infof("Initializing client")
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}
	klog.Infof("Client initialized")
	k8.Client = client
	return nil
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	var configPath string

	klog.Infof("Loading config")

	configPath = os.Getenv("CONTROLLER_CONFIG")
	if !(len(configPath) > 1) {
		configPath = DEFAULT_CONFIG_PATH
	}
	klog.Infof("Config Path - %s", configPath)
	config, err := loadConfigFromPath(configPath)
	if err != nil {
		klog.Fatal(err)
	}

	_conf, _ := json.MarshalIndent(config, "", "\t")
	klog.Infof("Config is %s", string(_conf))

	m := new(ebpfModule)
	m.config = config

	host := os.Getenv(MY_HOST)
	if !(len(host) > 1) {
		klog.Fatal("Could not retrieve host information")
	}

	k8client := new(KubeClient)
	err = k8client.InitKubeClient()
	if err != nil {
		klog.Fatal(err)
	}

	c := NewPodNetworkController(k8client.Client, m, host)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	stop := make(chan struct{})
	defer close(stop)

	defer runtime.HandleCrash()

	err = c.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}

	select {
	case <-sigChannel:
		klog.Infof("Received SIGTERM. Exiting")
		os.Exit(0)
	}
}
