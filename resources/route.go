package resources

import (
	"context"
	"fmt"
	"path"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/dshelley66/gcp-nuke/pkg/gcputil"
	"github.com/dshelley66/gcp-nuke/pkg/types"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

const ResourceTypeRoute = "Route"

type Route struct {
	name         string
	network      string
	creationDate string
	operation    *compute.Operation
}

func init() {
	register(ResourceTypeRoute, GetRouteClient, ListRoutes)
}

func GetRouteClient(project *gcputil.Project) (gcputil.GCPClient, error) {
	if client, ok := project.GetClient(ResourceTypeRoute); ok {
		return client, nil
	}

	client, err := compute.NewRoutesRESTClient(project.GetContext(), project.Creds.GetNewClientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("failed to create routes client: %v", err)
	}
	project.AddClient(ResourceTypeRoute, client)
	return client, nil
}

func ListRoutes(project *gcputil.Project, client gcputil.GCPClient) ([]Resource, error) {
	firewallsClient := client.(*compute.RoutesClient)

	resources := make([]Resource, 0)
	req := &computepb.ListRoutesRequest{
		Project: project.Name,
		Filter:  &noDefaultNetworkFilter,
	}

	it := firewallsClient.List(project.GetContext(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list routes: %v", err)
		}
		resources = append(resources, &Route{
			name:         *resp.Name,
			network:      *resp.Network,
			creationDate: *resp.CreationTimestamp,
		})
	}
	return resources, nil
}

func (x *Route) Remove(project *gcputil.Project, client gcputil.GCPClient) (err error) {
	firewallsClient := client.(*compute.RoutesClient)

	req := &computepb.DeleteRouteRequest{
		Route:   x.name,
		Project: project.Name,
	}

	x.operation, err = firewallsClient.Delete(project.GetContext(), req)
	if err != nil {
		return err
	}

	return nil
}

func (x *Route) GetOperationError(ctx context.Context) error {
	return getComputeOperationError(ctx, x.operation)
}

func (x *Route) String() string {
	return x.name
}

func (x *Route) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", x.name)
	properties.Set("Network", path.Base(x.network))
	properties.Set("CreationDate", x.creationDate)

	return properties
}
