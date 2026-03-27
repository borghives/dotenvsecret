package dotenvsecret

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type GCPSecretManager struct{}

func NewGCPSecretManager() *GCPSecretManager {
	return &GCPSecretManager{}
}

func getProjectNumber(ctx context.Context, projectID string) (string, error) {
	client, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create projects client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s", projectID)
	req := &resourcemanagerpb.GetProjectRequest{
		Name: name,
	}
	project, err := client.GetProject(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get project: %w", err)
	}
	parts := strings.Split(project.Name, "/")
	if len(parts) == 0 {
		return "", errors.New("invalid project name returned")
	}
	return parts[len(parts)-1], nil
}

func (m *GCPSecretManager) AccessSecret(ctx context.Context, secretID, versionID string) (string, error) {
	projectNum := os.Getenv("GOOGLE_CLOUD_PROJECT_NUM")
	if projectNum == "" {
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			projectID = os.Getenv("PROJECT_ID")
		}
		if projectID == "" {
			return "", errors.New("Project ID is missing. Set GOOGLE_CLOUD_PROJECT or PROJECT_ID environment variable.")
		}
		var err error
		projectNum, err = getProjectNumber(ctx, projectID)
		if err != nil {
			return "", err
		}
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectNum, secretID, versionID)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}
