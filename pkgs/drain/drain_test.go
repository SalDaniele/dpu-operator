package drain

import (
	"context"
	"testing"

	"github.com/openshift/dpu-operator/internal/testutils"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestKindCluster(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Kind Cluster Test Suite")
}

var _ = ginkgo.Describe("Kind Cluster Setup and Validation", ginkgo.Ordered, func() {
	var restConfig *rest.Config
	var k8sClient client.Client
	var kindCluster testutils.KindCluster

	// Setup for the whole test suite
	ginkgo.BeforeAll(func() {
		opts := zap.Options{
			Development: true,
		}
		ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		// Initialize KindCluster
		kindCluster = testutils.KindCluster{Name: "test-cluster"}

		// Create the Kind cluster
		restConfig = kindCluster.EnsureExists()

		// Create a Kubernetes client from the rest config
		var err error
		k8sClient, err = client.New(restConfig, client.Options{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})
	// Cleanup after all tests are complete
	ginkgo.AfterAll(func() {
		// Clean up the Kind cluster after the test
		kindCluster.EnsureDeleted()
	})
	// Test the node readiness in a separate context
	ginkgo.Context("nodes should be running", func() {
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
