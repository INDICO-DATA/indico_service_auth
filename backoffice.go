package indicoserviceauth

import (
	"context"
	"fmt"

	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	resourcesClient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceAccountClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (client *Client) CreateServiceAccount(ctx context.Context, displayName string, name string, description string) (*serviceAccountClient.CreateServiceAccountResponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createServiceAccount := &serviceAccountClient.CreateServiceAccountRequest{
		DisplayName: displayName,
		Name:        name,
		Description: description,
	}

	return client.serviceAccountService.CreateServiceAccount(ctx, createServiceAccount)
}

func (client *Client) CreateResource(ctx context.Context, name string, description string) (*resourcesClient.Resource, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createResourceRequest := &resourcesClient.CreateResourceRequest{
		Name:        name,
		Description: description,
	}

	return client.resourcesService.CreateResource(ctx, createResourceRequest)
}

func (client *Client) CreateResourceScope(ctx context.Context, label string, name string, description string, resourceID int32) (*resourcesClient.ResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createResourceScopeRequest := &resourcesClient.CreateResourceScopeRequest{
		Label:       label,
		Name:        name,
		Description: description,
		ResourceId:  resourceID,
	}

	return client.resourcesService.CreateResourceScope(ctx, createResourceScopeRequest)
}

func (client *Client) ListResources(ctx context.Context, principal string) (*resourcesClient.ListResource, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.resourcesService.ListResources(ctx, &emptypb.Empty{})
}

func (client *Client) GetResourceScope(ctx context.Context, resourceID int32) (*resourcesClient.ListResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	getResourceScopeRequest := &resourcesClient.GetResourceScopeRequest{
		ResourceId: resourceID,
	}

	return client.resourcesService.GetResourceScope(ctx, getResourceScopeRequest)
}

func (client *Client) CreateClient(ctx context.Context, principal string, clientType string) (*clientsClient.CreateClientResponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createClientRequest := &clientsClient.CreateClientRequest{
		Principal: principal,
		Type:      clientType,
	}

	return client.clientsService.CreateClient(ctx, createClientRequest)
}
