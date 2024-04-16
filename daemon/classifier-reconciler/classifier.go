package classifierreconciler

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configv1 "github.com/openshift/dpu-operator/api/v1"
)

// ClassifierConfigReconciler reconciles a ClassifierConfig object
type ClassifierConfigReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=config.openshift.io,resources=classifierconfigs,verbs=get;list;watch;create;update;patch;delet
//+kubebuilder:rbac:groups=config.openshift.io,resources=classifierconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.openshift.io,resources=classifierconfigs/finalizers,verbs=update
func (r *ClassifierConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    logger.Info("ClassifierConfigReconciler")

    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClassifierConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&configv1.ClassifierConfig{}).
        Complete(r)
}
