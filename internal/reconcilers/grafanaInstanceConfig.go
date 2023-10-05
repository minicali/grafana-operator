package reconcilers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"encoding/json"
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

	// Converting the json format to golang to map[string]interface{}
	parsedINI, err := convertToInterfaceMap(cr.Spec.INIConfig)
	if err != nil {
		log.Error(err, "Failed to convert INI input")
		return err
	}

	// Consolidate the config
	consolidatedConfig := consolidateConfig(parsedINI)

	// Get the INI content from the CR's spec.iniConfig
	iniContent, err := helpers.ToINIConfig(consolidatedConfig)
	if err != nil {
		log.Error(err, "Failed to generate INI content")
		return err
	}

	// Define a new ConfigMap object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helpers.GetPrefixedName(cr.Name, "config-ini"),
			Namespace: cr.Namespace,
			Labels:    helpers.GetGrafanaLabels(cr.Name, "config"),
		},
		Data: map[string]string{
			"grafana.ini": iniContent,
		},
	}

	// Check if this ConfigMap already exists
	found := &corev1.ConfigMap{}
	if err = r.Client.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: cr.Namespace}, found); err != nil && errors.IsNotFound(err) {
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
		// ConfigMap already exists, check if it needs to be updated
		if !reflect.DeepEqual(configMap.Data, found.Data) {
			// Update the found object and write the result back if there are any changes
			found.Data = configMap.Data
			log.Info("Updating ConfigMap")
			err = r.Client.Update(ctx, found)
			if err != nil {
				log.Error(err, "Failed to update ConfigMap")
				return err
			}
		} else {
			log.Info("Skip reconcile: ConfigMap already exists and is up-to-date")
		}
	}

	return nil
}

func convertToInterfaceMap(input map[string]apiextensionsv1.JSON) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})

	for key, value := range input {
		var tmp map[string]interface{}
		if err := json.Unmarshal(value.Raw, &tmp); err != nil {
			return nil, err
		}
		parsed[key] = tmp
	}

	return parsed, nil
}

func consolidateConfig(userConfig map[string]interface{}) map[string]interface{} {

	defaults := map[string]interface{}{
		"metrics": map[string]interface{}{
			"enabled":             true,
			"disable_total_stats": false,
		},
	}

	return helpers.MergeMaps(defaults, userConfig)
}
