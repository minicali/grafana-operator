package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PVCReconciler struct {
	Client client.Client
}

func NewPVCReconciler(client client.Client) *PVCReconciler {
	return &PVCReconciler{
		Client: client,
	}
}

func (r *PVCReconciler) Reconcile(ctx context.Context, cr *v1alpha1.GrafanaInstance, log logr.Logger) error {
	log = log.WithValues("Resource", "PVC")
	log.Info("Reconciling PVC")

	// Define a new PVC object
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helpers.GetPrefixedName(cr.Name, "pvc"),
			Namespace: cr.Namespace,
			Labels:    helpers.GetGrafanaLabels(cr.Name, "pvc"),
		},
		Spec: getGrafanaPvcSpec(),
	}

	// Check if this PVC already exists
	found := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: pvc.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the PVC since it doesn't exist
		log.Info("Creating a new PVC")
		err = r.Client.Create(ctx, pvc)
		if err != nil {
			log.Error(err, "Failed to create PVC")
			return err
		}
	} else if err != nil {
		log.Error(err, "Failed to get PVC")
		return err
	} else {
		log.Info("Skip reconcile: PVC already exists")
	}

	return nil
}

func getGrafanaPvcSpec() corev1.PersistentVolumeClaimSpec {
	return corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
		},
	}
}
