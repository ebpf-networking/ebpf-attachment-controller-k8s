package main

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	informersCoreV1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	klog "k8s.io/klog/v2"
)

const (
	EBPF_ATTACHMENT        string = "ebpf-attachment"
	EBPF_ATTACHMENT_STATUS string = "ebpf-attachment-status"
)

const (
	EBPF_ATTACHED string = "attached"
	EBPF_FAILED   string = "failed"
)

type PodNetworkController struct {
	clientset       kubernetes.Interface
	informerFactory informers.SharedInformerFactory
	podInformer     informersCoreV1.PodInformer
	ebpfController  *ebpfModule
	host            string
}

func (c *PodNetworkController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far.
	klog.Infof("Controller start")

	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *PodNetworkController) isOurPod(pod *v1.Pod) bool {
	// If this replica of daemonset is running on the same host as the pod
	// then we take the ownership else someone else takes the ownership
	return pod.Spec.NodeName == c.host
}

func (c *PodNetworkController) hasEBPFAttachment(pod *v1.Pod) bool {
	if !apimachinery.HasAnnotation(pod.ObjectMeta, EBPF_ATTACHMENT) {
		klog.Infof("Pod %s doesnt contain %s annotation", pod.Name, EBPF_ATTACHMENT)
		return false
	}
	return true
}

func (c *PodNetworkController) isEBPFAttached(pod *v1.Pod) bool {
	if !apimachinery.HasAnnotation(pod.ObjectMeta, EBPF_ATTACHMENT_STATUS) {
		klog.Infof("Pod %s doesnt contain %s annotation", pod.Name, EBPF_ATTACHMENT_STATUS)
		return false
	}
	annotations := pod.GetAnnotations()
	status := annotations[EBPF_ATTACHMENT_STATUS]

	if (status != EBPF_ATTACHED) && (status != EBPF_FAILED) {
		return false
	}
	klog.Infof("Pod %s contains annotation %s - %s",
		pod.Name, EBPF_ATTACHMENT_STATUS, status)
	return true
}

func (c *PodNetworkController) setAttachmentStatus(pod *v1.Pod, status string) error {
	annotations := pod.ObjectMeta.GetAnnotations()
	annotations[EBPF_ATTACHMENT_STATUS] = status

	copyPod := pod.DeepCopy()
	copyPod.SetAnnotations(annotations)

	_, err := c.clientset.CoreV1().Pods(pod.Namespace).Update(context.TODO(), copyPod, apimachinery.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Failed to update pod: %s", err)
	}

	klog.Infof("Pod attachment status updated")
	return nil
}

func (c *PodNetworkController) podAdd(_pod interface{}) {
	pod := _pod.(*v1.Pod)
	if !c.isOurPod(pod) {
		return
	}
	klog.Infof("POD CREATED: %s/%s", pod.Namespace, pod.Name)
}

func (c *PodNetworkController) podUpdate(_oldPod, _newPod interface{}) {
	oldPod := _oldPod.(*v1.Pod)
	newPod := _newPod.(*v1.Pod)

	if !c.isOurPod(newPod) {
		return
	}

	klog.Infof(
		"POD UPDATED. old %s/%s %s to %s/%s %s",
		oldPod.Namespace, oldPod.Name, oldPod.Status.Phase,
		newPod.Namespace, newPod.Name, newPod.Status.Phase,
	)
	if newPod.Status.Phase == v1.PodRunning {
		if c.hasEBPFAttachment(newPod) && !c.isEBPFAttached(newPod) {
			var status string
			klog.Infof("POD is in Running phase. %s/%s", newPod.Namespace, newPod.Name)
			annotations := newPod.ObjectMeta.GetAnnotations()
			program := annotations[EBPF_ATTACHMENT]
			err := c.ebpfController.AttachEBPFNetwork(newPod, program)
			if err != nil {
				klog.Info("ebpf attachment failed. %s", err)
				status = EBPF_FAILED
			}
			klog.Info("ebpf attachment done.")
			status = EBPF_ATTACHED
			c.setAttachmentStatus(newPod, status)
		}
	}
}

func (c *PodNetworkController) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	if !c.isOurPod(pod) {
		return
	}

	klog.Infof("POD DELETED: %s/%s", pod.Namespace, pod.Name)
	klog.Infof("Removing bpf attachment")
	if c.hasEBPFAttachment(pod) {
		annotations := pod.ObjectMeta.GetAnnotations()
		program := annotations[EBPF_ATTACHMENT]
		c.ebpfController.DeleteEBPFNetwork(pod, program)
	}
}

func NewPodNetworkController(clientset kubernetes.Interface, ebpfController *ebpfModule, host string) *PodNetworkController {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Hour*24)
	podInformer := informerFactory.Core().V1().Pods()

	c := new(PodNetworkController)
	c.clientset = clientset
	c.informerFactory = informerFactory
	c.podInformer = podInformer
	c.ebpfController = ebpfController
	c.host = host

	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}
