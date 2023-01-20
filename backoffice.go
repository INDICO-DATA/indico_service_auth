package indicoserviceauth

import (
	"context"
	"fmt"

	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	resourcesClient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceAccountClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ResourceScope struct {
	ResourceScopesID int    `json:"resource_scopes_id"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	Description      string `json:"description,omitempty"`
	ResourceID       int    `json:"resource_id"`
}

type ResourceScopeList struct {
	Data []*ResourceScope `json:"data"`
}

func (rsl *ResourceScopeList) createResourceScopeRequestMapper() *resourcesClient.CreateResourceScopeRequest {
	data := new(resourcesClient.CreateResourceScopeRequest)

	for _, item := range rsl.Data {
		data.Data = append(data.Data, &resourcesClient.ResourceScope{
			Name:        item.Name,
			Label:       item.Label,
			Description: item.Description,
			ResourceId:  int32(item.ResourceID),
		})
	}

	return data
}

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

func (client *Client) CreateResourceScope(ctx context.Context, payload *ResourceScopeList) (*resourcesClient.CreateResourceScopeRequest, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.resourcesService.CreateResourceScope(ctx, payload.createResourceScopeRequestMapper())

}

func (client *Client) ListResources(ctx context.Context) (*resourcesClient.ListResource, error) {
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

func (client *Client) CreateClient(ctx context.Context, principal string, clientType string) (*clientsClient.Client, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createClientRequest := &clientsClient.CreateClientRequest{
		Principal: principal,
		Type:      clientType,
	}

	return client.clientsService.CreateClient(ctx, createClientRequest)
}

func (client *Client) CreateRole(ctx context.Context, clientID int32, resourceScopeID int32) (*clientsClient.Role, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	clientRole := &clientsClient.CreateRoleRequest{
		ClientId:        clientID,
		ResourceScopeId: resourceScopeID,
	}

	return client.clientsService.CreateRole(ctx, clientRole)
}

func (client *Client) DeleteClient(ctx context.Context, clientID int32) (*clientsClient.DeleteClientReponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	deleteClientRequest := &clientsClient.DeleteClientRequest{
		ClientId: clientID,
	}

	return client.clientsService.DeleteClient(ctx, deleteClientRequest)
}

func (client *Client) DeleteRole(ctx context.Context, roleID int32) (*clientsClient.DeleteRoleReponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	deleteRoleRequest := &clientsClient.DeleteRoleRequest{
		RoleId: roleID,
	}

	return client.clientsService.DeleteRole(ctx, deleteRoleRequest)
}

func (client *Client) ListClients(ctx context.Context) (*clientsClient.ListClientResponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.clientsService.ListClients(ctx, &emptypb.Empty{})
}

func (client *Client) ListClientScopes(ctx context.Context, clientID int32) (*resourcesClient.ListResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	listClientScopesRequest := &resourcesClient.ListClientScopesRequest{
		ClientId: clientID,
	}

	return client.resourcesService.ListClientScopes(ctx, listClientScopesRequest)
}
