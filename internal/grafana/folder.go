package grafana

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/minicali/grafana-operator/api/v1alpha1"
)

// getFolderUIDByName retrieves the UID of a Grafana folder by its name.
// It returns the UID and a error indicating whether the folder was found.
func (gc *GrafanaClient) getFolderUIDByName(log logr.Logger, folderName string) (string, error) {
	log = log.WithValues("Resource", "Folder")

	// exception, pre-existing folder aren't returned from API
	if IsGeneralFolder(folderName) {
		return "0", nil
	}

	log.Info("Listing Grafana folders")
	// Fetch the list of folders from Grafana
	folders, err := gc.Client.Folders()
	if err != nil {
		log.Error(err, "Failed to list Grafana folders")
		return "", err
	}

	// Loop through the folders to find the one that matches `cr.Spec.Folder`
	for _, folder := range folders {
		if strings.EqualFold(folder.Title, folderName) {
			log.Info("Found matching Grafana folder", "folderUID", folder.UID)
			return folder.UID, nil
		}
	}

	log.Info("No matching Grafana folder found", "folderName", folderName)
	return "", fmt.Errorf("Folder '%s' not found", folderName)
}

// EnsureFolder ensures that a Grafana folder exists.
// If the GrafanaDashboard's Status includes a folder ID, it updates the folder with the name.
// Otherwise, it creates a new folder and returns its ID.
func (c *GrafanaClient) EnsureFolder(log logr.Logger, cr *grafanav1alpha1.GrafanaDashboard) (string, error) {
	// Check if UID exists in the status
	existingUID := cr.Status.FolderUID

	// General folder already exist
	if IsGeneralFolder(cr.Spec.Folder) {
		return "", nil
	}

	// If UID exists, update the folder
	if existingUID != "" {
		log.Info("Updating existing Grafana folder", "UID", existingUID)
		err := c.Client.UpdateFolder(existingUID, cr.Spec.Folder)
		if err != nil {
			return "", fmt.Errorf("failed to update Grafana folder: %w", err)
		}
		return existingUID, nil
	}

	// Otherwise, create a new folder
	log.Info("Creating new Grafana folder", "Title", cr.Spec.Folder)
	resp, err := c.Client.NewFolder(cr.Spec.Folder)
	if err != nil {
		return "", fmt.Errorf("failed to create new Grafana folder: %w", err)
	}

	return resp.UID, nil
}

// GetFolderIDByUID retrieves the folder ID based on its UID.
// Returns an error if UID is empty or if the API call fails.
func (gc *GrafanaClient) GetFolderIDByUID(uid string) (int64, error) {
	if uid == "" {
		return -1, errors.New("UID cannot be empty")
	}

	folder, err := gc.Client.FolderByUID(uid)
	if err != nil {
		return -1, fmt.Errorf("failed to fetch folder by UID: %w", err)
	}

	return folder.ID, nil
}

// IsGeneralFolder checks if the folder is the pre-existing one General.
// This folder is an exception and isn't included in the API.
func IsGeneralFolder(folderName string) bool {
	if strings.EqualFold(folderName, grafanav1alpha1.GrafanaGeneralFolder) {
		return true
	}
	return false
}
