package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceReconciler struct {
	Client client.Client
}

func NewServiceReconciler(client client.Client) *ServiceReconciler {
	return &ServiceReconciler{
		Client: client,
	}
}

func (r *ServiceReconciler) Reconcile(ctx context.Context, cr *grafanav1alpha1.GrafanaInstance, log logr.Logger) error {
	log.Info("Reconciling Service")

	// Define a new Service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helpers.GetPrefixedName(cr.Name, "service"),
			Namespace: cr.Namespace,
			Labels:    helpers.GetGrafanaLabels(cr.Name, "service"),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       3000,
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("http-grafana"),
				},
			},
			Selector: helpers.GetGrafanaLabels(cr.Name, "deployment"),
		},
	}

	// Check if this Service already exists
	found := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the Service since it doesn't exist
		log.Info("Creating a new Service")
		err = r.Client.Create(ctx, service)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		log.Info("Skip reconcile: Service already exists")
	}

	return nil
}
