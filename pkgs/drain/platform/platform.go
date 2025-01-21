package platforms

// Dummy implementation of PlatformHelper and NewHelpers
import "k8s.io/client-go/kubernetes"

type Interface interface {
    OpenshiftBeforeDrainNode(ctx context.Context, node *corev1.Node) (bool, error)
    OpenshiftAfterCompleteDrainNode(ctx context.Context, node *corev1.Node) (bool, error)
}

type PlatformHelper struct{}

func (p *PlatformHelper) OpenshiftBeforeDrainNode(ctx context.Context, node *corev1.Node) (bool, error) {
    return true, nil
}

func (p *PlatformHelper) OpenshiftAfterCompleteDrainNode(ctx context.Context, node *corev1.Node) (bool, error) {
    return true, nil
}

// NewHelpers creates and returns a PlatformHelper
func NewHelpers() Interface {
    return &PlatformHelper{}
}

