package grafana

import (
	"fmt"

	"github.com/go-logr/logr"
	grapi "github.com/grafana/grafana-api-golang-client"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
)

func (gc *GrafanaClient) UpsertDashboard(log logr.Logger, cr *grafanav1alpha1.GrafanaDashboard, folderID int64) error {
	log = log.WithValues("Resource", "Dashboard")
	log.Info("Reconciling GrafanaDashboard")

	// TODO
	dashboardModel := map[string]interface{}{}

	resp, err := gc.Client.NewDashboard(grapi.Dashboard{
		Model:     dashboardModel,
		FolderID:  folderID,
		Overwrite: true,
		Message:   "",
	})

	if err != nil {
		log.Error(err, "Failed to create/update Grafana dashboard")
		return err
	}

	if resp.Status != "success" {
		log.Error(nil, "Error creating dashboard, status was not 'success'", "status", resp.Status)
		return fmt.Errorf("error creating dashboard, status was %v", resp.Status)
	}

	log.Info("Successfully created/updated Grafana dashboard", "dashboardName", cr.Spec.Name)
	return nil
}
