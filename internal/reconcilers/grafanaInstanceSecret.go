package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"crypto/sha256"
	"encoding/hex"
)

type SecretReconciler struct {
	Client client.Client
}

func NewSecretReconciler(client client.Client) *SecretReconciler {
	return &SecretReconciler{
		Client: client,
	}
}

func (r *SecretReconciler) Reconcile(ctx context.Context, cr *v1alpha1.GrafanaInstance, log logr.Logger) error {
	log = log.WithValues("Resource", "Secret")
	log.Info("Reconciling Secret")

	// Generate credentials
	username := helpers.GenerateRandomString(10)
	password := helpers.GenerateRandomString(10)

	// Define a new Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.CredentialsSecretName,
			Namespace: cr.Namespace,
			Annotations: map[string]string{
				"generated-by": "grafana-operator", // New annotation
			},
		},
		StringData: map[string]string{
			"admin_username": username,
			"admin_password": password,
		},
		Type: corev1.SecretTypeOpaque,
	}

	// Generate a hash of the secret data
	hash := sha256.Sum256([]byte(username + password))
	hashStr := hex.EncodeToString(hash[:])

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: secret.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the Secret since it doesn't exist
		log.Info("Creating a new Secret")
		secret.Annotations["checksum"] = hashStr

		err = r.Client.Create(ctx, secret)
		if err != nil {
			log.Error(err, "Failed to create Secret")
			return err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Secret")
		return err
	} else {
		// Check if hash has changed
		oldHash, ok := found.Annotations["checksum"]
		if !ok || oldHash != hashStr {
			log.Info("Secret has changed, updating Secret")

			// Update annotation on Secret
			found.Annotations["checksum"] = hashStr
			if err := r.Client.Update(ctx, found); err != nil {
				log.Error(err, "Failed to update Secret")
				return err
			}
		} else {
			log.Info("Skip reconcile: Secret already exists")
		}
	}

	return nil
}
