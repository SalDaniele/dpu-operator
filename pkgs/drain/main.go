package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// set up args, namely kubeconfig / node name
	fmt.Printf("starting main program\n")
	kubeconfig := flag.String("kubeconfig", "test", "Location of your kubeconfig")
	nodeName := flag.String("node", "", "Name of the node to drain.")
	flag.Parse()
	fmt.Printf("flag: %v:\n", *kubeconfig)
	fmt.Printf("node: %v:\n", *nodeName)

	// Load Kubernetes config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	node, err := kubeClient.CoreV1().Nodes().Get(context.TODO(), *nodeName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error retrieving node: %v\n", err)
		os.Exit(1)
	}

	// Print out the node details
	fmt.Printf("Node Name: %s\n", node.Name)
	fmt.Printf("Node Status: %s\n", node.Status.Phase)
	fmt.Printf("Node Labels: %v\n", node.Labels)


	// Now do the platform specific stuff



}
