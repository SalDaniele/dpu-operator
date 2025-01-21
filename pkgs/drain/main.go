package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	platforms "github.com/SalDaniele/dpu-operator/pkgs/drain/platform"
	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/drain"
)

func main() {
	// Load kubeconfig from file or default location
	kubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = fmt.Sprintf("%s/.kube/config", home)
	}
	kubeconfigFlag := flag.String("kubeconfig", kubeconfig, "path to the kubeconfig file")
	flag.Parse()

	// Create Kubernetes client from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigFlag)
	if err != nil {
		log.Fatalf("failed to build kubeconfig: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kube client: %v", err)
	}

	// Create the Drainer object
	platformHelpers := platforms.NewHelpers() // Implement your platform helpers if needed
	drainer, err := drain.NewDrainer(platformHelpers)
	if err != nil {
		log.Fatalf("failed to create drainer: %v", err)
	}

	// Get the node you want to drain
	nodeName := "your-node-name" // Set the node name
	node, err := kubeClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("failed to get node %s: %v", nodeName, err)
	}

	// Drain the node
	drainSuccess, err := drainer.DrainNode(context.TODO(), node, true)
	if err != nil {
		log.Fatalf("failed to drain node %s: %v", nodeName, err)
	}
	if drainSuccess {
		log.Printf("Node %s drained successfully", nodeName)
	} else {
		log.Printf("Failed to drain node %s, retrying...", nodeName)
	}

	// Complete the drain
	completeSuccess, err := drainer.CompleteDrainNode(context.TODO(), node)
	if err != nil {
		log.Fatalf("failed to complete drain on node %s: %v", nodeName, err)
	}
	if completeSuccess {
		log.Printf("Drain on node %s completed successfully", nodeName)
	} else {
		log.Printf("Failed to complete drain on node %s", nodeName)
	}
}
