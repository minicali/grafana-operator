package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
)

type StageReconciler interface {
	Reconcile(ctx context.Context, cr *grafanav1alpha1.GrafanaInstance, log logr.Logger) error
}
