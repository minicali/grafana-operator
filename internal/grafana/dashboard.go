package grafana

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	grapi "github.com/grafana/grafana-api-golang-client"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
	"github.com/minicali/grafana-operator/internal/helpers"
)

func (gc *GrafanaClient) UpsertDashboard(log logr.Logger, cr *grafanav1alpha1.GrafanaDashboard, folderUID string) (string, error) {
	log = log.WithValues("Resource", "Dashboard")

	dashboardModel, err := getModelFromCR(cr)
	if err != nil {
		log.Error(err, "Failed to create/update Grafana dashboard")
		return "", err
	}

	resp, err := gc.Client.NewDashboard(grapi.Dashboard{
		Model:     dashboardModel,
		FolderUID: folderUID,
		Overwrite: true,
		Message:   "Upserted by Grafana-Operator at " + time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		log.Error(err, "Failed to create/update Grafana dashboard")
		return "", err
	}

	if resp.Status != "success" {
		log.Error(nil, "Error creating dashboard, status was not 'success'", "status", resp.Status)
		return "", fmt.Errorf("error creating dashboard, status was %v", resp.Status)
	}

	log.Info("Successfully created/updated Grafana dashboard", "dashboardName", cr.Spec.Name)
	return resp.UID, nil
}

func (gc *GrafanaClient) DeleteDashboard(log logr.Logger, dashboardUID string) error {
	log = log.WithValues("Resource", "Dashboard")

	if dashboardUID == "" {
		log.Error(nil, "Error deleting dashboard, UID is missing")
		return fmt.Errorf("error deleting dashboard, UID is missing")
	}

	err := gc.Client.DeleteDashboardByUID(dashboardUID)
	if err != nil {
		log.Error(err, "Failed to delete Grafana dashboard")
		return err
	}

	log.Info("Successfully deleted Grafana dashboard", "dashboardUID", dashboardUID)
	return nil
}

func getModelFromCR(cr *grafanav1alpha1.GrafanaDashboard) (map[string]interface{}, error) {
	switch {
	case cr.Spec.Json.Raw != nil:
		return helpers.UnmarshalJSONToMap(cr.Spec.Json)
	default:

		return nil, errors.New("no valid dashboard model found")
	}
}
