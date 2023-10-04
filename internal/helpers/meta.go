package helpers

import "fmt"

// GetGrafanaLabels returns a map of standard labels for a Grafana instance
func GetGrafanaLabels(instanceName string, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "grafana-dashboard",
		"app.kubernetes.io/instance":   instanceName,
		"app.kubernetes.io/version":    "10.0.0",
		"app.kubernetes.io/managed-by": "grafana-operator",
		"app.kubernetes.io/component":  component,
	}
}

// GetPrefixedName returns a Kubernetes resource name with the given suffix and custom resource name
func GetPrefixedName(crName string, suffix string) string {
	return fmt.Sprintf("%s-%s", crName, suffix)
}
