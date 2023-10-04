package grafana

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// IsGrafanaHealthy checks the health of the Grafana instance using grafanaClient.Health()
func (gc *GrafanaClient) IsGrafanaHealthy() (bool, error) {
	healthStatus, err := gc.Client.Health()
	if err != nil {
		return false, fmt.Errorf("error checking Grafana health: %v", err)
	}

	if healthStatus.Database != "ok" || healthStatus.Version == "" {
		return false, fmt.Errorf("Grafana is not healthy: Database: %s, Version: %s", healthStatus.Database, healthStatus.Version)
	}
	return true, nil

	// return false, fmt.Errorf("received empty health status")
}

func CreateGrafanaClientFromSecret(ctx context.Context, grafanaInstanceSecret *corev1.Secret, grafanaURL string) (*GrafanaClient, error) {
	// Retrieve username and password from the secret
	username := string(grafanaInstanceSecret.Data["admin_username"])
	password := string(grafanaInstanceSecret.Data["admin_password"])

	// Create Grafana client
	auth := &BasicAuthenticator{
		Username: username,
		Password: password,
	}
	timeout := 5 // adjust as needed

	grafanaClient, err := NewClient(grafanaURL, time.Duration(timeout)*time.Second, auth)
	if err != nil {
		return nil, err
	}

	return grafanaClient, nil
}
