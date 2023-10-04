package reconcilers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeploymentReconciler struct {
	Client client.Client
}

func NewDeploymentReconciler(client client.Client) *DeploymentReconciler {
	return &DeploymentReconciler{
		Client: client,
	}
}
func (r *DeploymentReconciler) Reconcile(ctx context.Context, cr *v1alpha1.GrafanaInstance, log logr.Logger) error {
	log = log.WithValues("Resource", "Deployment")
	log.Info("Reconciling Deployment")

	// Define a new Deployment object
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helpers.GetPrefixedName(cr.Name, "ui"),
			Namespace: cr.Namespace,
			Labels:    helpers.GetGrafanaLabels(cr.Name, "deployment"),
		},
		Spec: getGrafanaDeploymentSpec(cr),
	}

	// Fetch the Secret to get the hash annotation
	secret := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: cr.Spec.CredentialsSecretName, Namespace: cr.Namespace}, secret)
	if err != nil {
		log.Error(err, "Failed to get Secret for hash annotation")
		return err
	}

	hash, ok := secret.Annotations["checksum"]
	if !ok {
		log.Error(fmt.Errorf("checksum annotation not found"), "Checksum annotation not found on Secret", "SecretName", cr.Spec.CredentialsSecretName)
		return fmt.Errorf("checksum annotation not found on Secret %s", cr.Spec.CredentialsSecretName)
	}

	// Add or update the hash annotation on the Deployment
	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}
	deployment.Annotations["secret-checksum"] = hash

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Client.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the Deployment since it doesn't exist
		log.Info("Creating a new Deployment")
		err = r.Client.Create(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to create Deployment")
			return err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return err
	} else {
		log.Info("Skip reconcile: Deployment already exists")
	}

	return nil
}

func getGrafanaDeploymentSpec(cr *v1alpha1.GrafanaInstance) appsv1.DeploymentSpec {
	pvcName := helpers.GetPrefixedName(cr.Name, "pvc")

	return appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: helpers.GetGrafanaLabels(cr.Name, "deployment"),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: helpers.GetGrafanaLabels(cr.Name, "deployment"),
			},
			Spec: corev1.PodSpec{
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup: pointer.Int64Ptr(472),
					SupplementalGroups: []int64{
						0,
					},
				},
				Containers: []corev1.Container{
					{
						Name:            "grafana",
						Image:           cr.Spec.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 3000,
								Name:          "http-grafana",
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      pvcName,
								MountPath: "/var/lib/grafana",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name: "GF_SECURITY_ADMIN_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: cr.Spec.CredentialsSecretName,
										},
										Key: "admin_username",
									},
								},
							},
							{
								Name: "GF_SECURITY_ADMIN_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: cr.Spec.CredentialsSecretName,
										},
										Key: "admin_password",
									},
								},
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: pvcName,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: pvcName,
							},
						},
					},
				},
			},
		},
	}
}
