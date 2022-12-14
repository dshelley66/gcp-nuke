package resources

import (
	"context"
	"fmt"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/dshelley66/gcp-nuke/pkg/gcputil"
	"github.com/dshelley66/gcp-nuke/pkg/types"
	"google.golang.org/api/iterator"
)

const ResourceTypeSecret = "Secret"

type Secret struct {
	name         string
	labels       map[string]string
	creationDate string
}

func init() {
	register(ResourceTypeSecret, GetSecretClient, ListSecret)
}

func GetSecretClient(project *gcputil.Project) (gcputil.GCPClient, error) {
	if client, ok := project.GetClient(ResourceTypeSecret); ok {
		return client, nil
	}

	client, err := secretmanager.NewClient(project.GetContext(), project.Creds.GetNewClientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	project.AddClient(ResourceTypeSecret, client)
	return client, nil
}

func ListSecret(project *gcputil.Project, client gcputil.GCPClient) ([]Resource, error) {
	secretClient := client.(*secretmanager.Client)

	resources := make([]Resource, 0)
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", project.Name),
	}

	it := secretClient.ListSecrets(project.GetContext(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %v", err)
		}
		resources = append(resources, &Secret{
			name:         resp.Name,
			creationDate: resp.CreateTime.AsTime().Format(time.RFC3339),
			labels:       resp.GetLabels(),
		})
	}
	return resources, nil
}

func (x *Secret) Remove(project *gcputil.Project, client gcputil.GCPClient) error {
	secretClient := client.(*secretmanager.Client)

	req := &secretmanagerpb.DeleteSecretRequest{
		Name: x.name,
	}

	err := secretClient.DeleteSecret(project.GetContext(), req)
	if err != nil {
		return err
	}

	return nil
}

func (x *Secret) GetOperationError(_ context.Context) error {
	return nil
}

func (x *Secret) String() string {
	return x.name
}

func (x *Secret) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", x.name)
	properties.Set("CreationDate", x.creationDate)

	for labelKey, label := range x.labels {
		properties.SetTag(labelKey, label)
	}

	return properties
}
