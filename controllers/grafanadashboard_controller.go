/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/minicali/grafana-operator/internal/grafana"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
)

// GrafanaDashboardReconciler reconciles a GrafanaDashboard object
type GrafanaDashboardReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const grafanaDashboardFinalizer = "finalizer.grafana.minicali.com"

//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanadashboards/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanadashboards/finalizers,verbs=update
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanainstances,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaDashboard object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GrafanaDashboardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("grafanadashboard", req.NamespacedName)
	log.Info("Starting reconciliation")

	// your code to get the GrafanaDashboard resource
	grafanaDashboard := &grafanav1alpha1.GrafanaDashboard{}
	if err := r.Get(ctx, req.NamespacedName, grafanaDashboard); err != nil {
		log.Error(err, "unable to fetch GrafanaDashboard")
		return ctrl.Result{}, err
	}

	grafanaDashboard.SetDefaults()

	// Get the GrafanaInstance resource
	grafanaInstanceName := grafanaDashboard.Spec.GrafanaInstanceRef.Name
	grafanaInstanceNamespace := grafanaDashboard.Spec.GrafanaInstanceRef.Namespace

	grafanaInstance := &grafanav1alpha1.GrafanaInstance{}
	err := r.Get(ctx, client.ObjectKey{Name: grafanaInstanceName, Namespace: grafanaInstanceNamespace}, grafanaInstance)
	if err != nil {
		log.Error(err, "unable to fetch GrafanaInstance")
		return ctrl.Result{}, err
	}

	// Get the service URL
	grafanaInstanceServiceURL := grafanaInstance.Status.GrafanaUI.ServiceURL
	if grafanaInstanceServiceURL == "" {
		log.Error(err, "failed to retrieve url")
		return ctrl.Result{}, err
	}
	// Get the secret
	grafanaInstanceSecret := &corev1.Secret{}
	err = r.Get(ctx, client.ObjectKey{Name: grafanaInstance.Spec.CredentialsSecretName, Namespace: grafanaInstance.Namespace}, grafanaInstanceSecret)
	if err != nil {
		log.Error(err, "unable to fetch Secret")
		return ctrl.Result{}, err
	}

	// Create Grafana client
	grafanaClient, err := grafana.CreateGrafanaClientFromSecret(ctx, grafanaInstanceSecret, grafanaInstanceServiceURL)
	if err != nil {
		log.Error(err, "unable to create Grafana Client")
		return ctrl.Result{}, err
	}

	// Check Grafana health
	_, err = grafanaClient.IsGrafanaHealthy()
	if err != nil {
		log.Error(err, "Failed to check Grafana health")
		return ctrl.Result{}, err
	}
	log.Info("Grafana connection healthy")

	// Check if the GrafanaDashboard resource is being deleted
	if grafanaDashboard.DeletionTimestamp != nil {
		// The object is being deleted
		if containsString(grafanaDashboard.ObjectMeta.Finalizers, grafanaDashboardFinalizer) {
			// Delete object from grafana here

			// Remove the finalizer from the list and update it.
			grafanaDashboard.ObjectMeta.Finalizers = removeString(grafanaDashboard.ObjectMeta.Finalizers, grafanaDashboardFinalizer)
			if err := r.Update(context.Background(), grafanaDashboard); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR, if it doesn't exist
	if !containsString(grafanaDashboard.ObjectMeta.Finalizers, grafanaDashboardFinalizer) {
		grafanaDashboard.ObjectMeta.Finalizers = append(grafanaDashboard.ObjectMeta.Finalizers, grafanaDashboardFinalizer)
		if err := r.Update(context.Background(), grafanaDashboard); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Requeue for periodic sync
	syncPeriod := grafanaDashboard.Spec.SyncPeriod.Duration
	return ctrl.Result{RequeueAfter: syncPeriod}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaDashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaDashboard{}).
		Complete(r)
}

// Utility functions
func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func removeString(slice []string, str string) []string {
	var newSlice []string
	for _, item := range slice {
		if item == str {
			continue
		}
		newSlice = append(newSlice, item)
	}
	return newSlice
}
