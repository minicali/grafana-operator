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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapReconciler struct {
	Client client.Client
}

func NewConfigMapReconciler(client client.Client) *ConfigMapReconciler {
	return &ConfigMapReconciler{
		Client: client,
	}
}

func (r *ConfigMapReconciler) Reconcile(ctx context.Context, cr *grafanav1alpha1.GrafanaInstance, log logr.Logger) error {
	log = log.WithValues("Resource", "ConfigMap")
	log.Info("Reconciling ConfigMap")

	// Define a new ConfigMap object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helpers.GetPrefixedName(cr.Name, "config-ini"),
			Namespace: cr.Namespace,
			Labels:    helpers.GetGrafanaLabels(cr.Name, "config"),
		},
		Data: map[string]string{
			"grafana.ini": `
[server]
  domain = sd
  root_url = https://%(domain)s
...
`,
		},
	}

	// Check if this ConfigMap already exists
	found := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the ConfigMap since it doesn't exist
		log.Info("Creating a new ConfigMap")
		err = r.Client.Create(ctx, configMap)
		if err != nil {
			log.Error(err, "Failed to create ConfigMap")
			return err
		}
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return err
	} else {
		log.Info("Skip reconcile: ConfigMap already exists")
	}

	return nil
}
