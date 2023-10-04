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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
	"github.com/minicali/grafana-operator/internal/reconcilers"
)

// GrafanaInstanceReconciler reconciles a GrafanaInstance object
type GrafanaInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type grafanaInstanceReconcileStages string

const (
	GrafanaInstanceStageConfigMap  grafanaInstanceReconcileStages = "config"
	GrafanaInstanceStagePVC        grafanaInstanceReconcileStages = "pvc"
	GrafanaInstanceStageSecret     grafanaInstanceReconcileStages = "secret"
	GrafanaInstanceStageDeployment grafanaInstanceReconcileStages = "deployment"
	GrafanaInstanceStageService    grafanaInstanceReconcileStages = "service"
)

var reconcileStages = []grafanaInstanceReconcileStages{
	GrafanaInstanceStageSecret,
	GrafanaInstanceStageConfigMap,
	GrafanaInstanceStagePVC,
	GrafanaInstanceStageDeployment,
	GrafanaInstanceStageService,
}

//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanainstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanainstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=grafana.minicali.com,resources=grafanainstances/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GrafanaInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithName("GrafanaInstanceController").WithValues("GrafanaInstance", req.NamespacedName)
	log.Info("Starting reconciliation")

	cr := &grafanav1alpha1.GrafanaInstance{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		log.Error(err, "Failed to get GrafanaInstance")
		return ctrl.Result{}, err
	}

	// Reconcile
	for _, stage := range reconcileStages {
		log.Info("Reconciling stage", "Stage", stage)
		stageReconciler := r.getReconcilerPerStage(stage)
		if err := stageReconciler.Reconcile(ctx, cr, log); err != nil {
			log.Error(err, "Failed to reconcile stage", "Stage", stage)
			return ctrl.Result{}, err
		}
	}

	// Status update
	// Fetch the Deployment
	deployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: helpers.GetPrefixedName(cr.Name, "ui"), Namespace: cr.Namespace}, deployment)
	if err != nil {
		log.Error(err, "Failed to get Deployment for status update")
		return ctrl.Result{}, err
	}

	// Update GrafanaInstanceStatus
	cr.Status.GrafanaUI.AvailableReplicas = fmt.Sprintf("%d/%d", deployment.Status.AvailableReplicas, *deployment.Spec.Replicas)
	cr.Status.GrafanaUI.Conditions = deployment.Status.Conditions

	// Update the GrafanaInstance status
	if err := r.Status().Update(ctx, cr); err != nil {
		log.Error(err, "Failed to update GrafanaInstance status")
		return ctrl.Result{}, err
	}

	log.Info("Finished reconciliation")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaInstance{}).
		Complete(r)
}

func (r *GrafanaInstanceReconciler) getReconcilerPerStage(stage grafanaInstanceReconcileStages) reconcilers.StageReconciler {
	switch stage {
	case GrafanaInstanceStageConfigMap:
		return reconcilers.NewConfigMapReconciler(r.Client)
	case GrafanaInstanceStagePVC:
		return reconcilers.NewPVCReconciler(r.Client)
	case GrafanaInstanceStageSecret:
		return reconcilers.NewSecretReconciler(r.Client)
	case GrafanaInstanceStageDeployment:
		return reconcilers.NewDeploymentReconciler(r.Client)
	case GrafanaInstanceStageService:
		return reconcilers.NewServiceReconciler(r.Client)
	default:
		return nil
	}
}
