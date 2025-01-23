package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/platforms"
	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/drain"
	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/vars"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/zap/zapcore"
)

	func main() {
	// Set up the logger as the very first thing
	log.SetLogger(zap.New(
		zap.WriteTo(GinkgoWriter),
		zap.Level(zapcore.Level(-2)),
		zap.UseDevMode(true)))
	// set up args, namely kubeconfig / node name
	fmt.Printf("starting main program\n")
	kubeconfig := flag.String("test-kubeconfig", "test", "Location of your kubeconfig")
	nodeName := flag.String("test-node", "", "Name of the node to drain.")
	flag.Parse()
	fmt.Printf("flag: %v:\n", *kubeconfig)
	fmt.Printf("node: %v:\n", *nodeName)

	// Load Kubernetes config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	vars.Config = config
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


	// Now do the drainer specific stuff
	fmt.Printf("creating platform helper\n")
	pH, err := platforms.NewDefaultPlatformHelper()
	if err != nil {
		fmt.Printf("couldn't create openshift context %v\n", err)
		os.Exit(1)
	}

	if pH == nil {
		fmt.Printf("nil platformHelpers")
		os.Exit(1)
	} else {
		fmt.Printf("not nil")
	}

	fmt.Printf("creating drainer helper\n")
	drainer, err := drain.NewDrainer(pH)
	if err != nil {
		fmt.Printf("Error creating Drainer: %v\n", err)
		os.Exit(1)
	}


	fmt.Printf("Draining node\n")
	drainSuccess, err := drainer.DrainNode(context.TODO(), node, true)
	if err != nil {
		fmt.Printf("Error draining node: %v\n", err)
		os.Exit(1)
	}

	if drainSuccess {
		fmt.Printf("Successfully drained node: %s\n", *nodeName)
	} else {
		fmt.Printf("Failed to drain node: %s\n", *nodeName)
	}
}
