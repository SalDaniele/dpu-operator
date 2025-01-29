package drain

import (
	"context"
	"os"
	"testing"

	"github.com/openshift/dpu-operator/internal/testutils"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/drain"
	"github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/platforms"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var restConfig *rest.Config
var k8sClient client.Client
var testCluster testutils.KindCluster

func createNewDrainer() (drain.DrainInterface, error) {
	platform, err := platforms.NewDefaultPlatformHelper()
	if err != nil {
		return nil, err
	}

	kclient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &drain.Drainer{
		kubeClient:      kclient,
		platformHelpers: platform,
	}, nil
}

func TestKindCluster(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Drain helper suite")
}

var _ = ginkgo.Describe("Kind Cluster Setup and Validation", ginkgo.Ordered, func() {

	// Setup for the whole test suite
	ginkgo.BeforeAll(func() {
		var err error

		// Setup logging
		opts := zap.Options{
			Development: true,
		}
		ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

		// Create a KindCluster
		testCluster = testutils.KindCluster{Name: "dpu-operator-test-cluster"}
		restConfig = testCluster.EnsureExists()

		// Create a Kubernetes client from the rest config
		k8sClient, err = client.New(restConfig, client.Options{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})
	ginkgo.AfterAll(func() {
		if os.Getenv("FAST_TEST") == "false" {
			testCluster.EnsureDeleted()
		}
	})
	// Test the node readiness in a separate context
	ginkgo.Context("When node drain has not been requested", func() {
		ginkgo.It("should have all nodes available", func() {
			ginkgo.By("checking that all nodes are ready and not draining")
			nodes := &corev1.NodeList{}
			err := k8sClient.List(context.Background(), nodes)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			for _, node := range nodes.Items {
				isReady := false
				for _, condition := range node.Status.Conditions {
					if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
						isReady = true
						break
					}
				}
				gomega.Expect(isReady).To(gomega.BeTrue(), "Node is ready: "+node.Name)

				// Check that the node is not unschedulable (draining)
				gomega.Expect(node.Spec.Unschedulable).To(gomega.BeFalse(), "Node is unschedulable (draining): "+node.Name)

				// Check that the node does not have a 'NoSchedule' taint (draining)
				draining := false
				for _, taint := range node.Spec.Taints {
					if taint.Effect == corev1.TaintEffectNoSchedule {
						draining = true
						break
					}
				}
				gomega.Expect(draining).To(gomega.BeFalse(), "Node is draining (has NoSchedule taint): "+node.Name)
			}
		})
		ginkgo.It("should have all nodes ready and not unschedulable", func() {
			ginkgo.By("checking that all nodes are ready and not draining")
			nodes := &corev1.NodeList{}
			err := k8sClient.List(context.Background(), nodes)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			for _, node := range nodes.Items {
				// Check that the node is ready
				isReady := false
				for _, condition := range node.Status.Conditions {
					if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
						isReady = true
						break
					}
				}
				gomega.Expect(isReady).To(gomega.BeTrue(), "Node is not ready: "+node.Name)

				// Check that the node is not unschedulable (draining)
				gomega.Expect(node.Spec.Unschedulable).To(gomega.BeFalse(), "Node is unschedulable (draining): "+node.Name)

				// Check that the node does not have a 'NoSchedule' taint (draining)
				draining := false
				for _, taint := range node.Spec.Taints {
					if taint.Effect == corev1.TaintEffectNoSchedule {
						draining = true
						break
					}
				}
				gomega.Expect(draining).To(gomega.BeFalse(), "Node is draining (has NoSchedule taint): "+node.Name)
			}
		})
	})
})
