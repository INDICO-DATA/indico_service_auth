package indicoserviceauth

import (
	"context"
	"fmt"
	"io"
	"log"

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

// Service Accounts
func (client *Client) CreateServiceAccount(ctx context.Context, displayName string, name string, description string) (*serviceAccountClient.ServiceAccount, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createServiceAccount := &serviceAccountClient.CreateRequest{
		DisplayName: displayName,
		Name:        name,
		Description: description,
	}

	return client.serviceAccountService.Create(ctx, createServiceAccount)
}

func (client *Client) ListServiceAccounts(ctx context.Context) ([]*serviceAccountClient.ServiceAccount, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	serviceAccounts := []*serviceAccountClient.ServiceAccount{}

	stream, err := client.serviceAccountService.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("error to stream list service account due to: %w", err)
	}

	for {
		serviceAccount, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming list service accounts: %w", err)
		}

		serviceAccounts = append(serviceAccounts, serviceAccount)
	}

	return serviceAccounts, nil
}

func (client *Client) GenerateCredentials(ctx context.Context, id string) (*serviceAccountClient.ServiceAccountCredentials, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.serviceAccountService.GenerateCredentials(ctx, &serviceAccountClient.CredentialsRequest{ServiceAccountId: id})
}

func (client *Client) DeleteServiceAccount(ctx context.Context, ids []string) ([]*serviceAccountClient.ServiceAccount, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.serviceAccountService.Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming delete service accounts: %w", err)
	}

	deletedServicesAccounts := []*serviceAccountClient.ServiceAccount{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming delete service accounts response: %s", err.Error())
			}

			deletedServicesAccounts = append(deletedServicesAccounts, in)
		}
	}()

	for _, accountID := range ids {
		if err := stream.Send(&serviceAccountClient.CredentialsRequest{ServiceAccountId: accountID}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming delete service accounts: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return deletedServicesAccounts, nil
}

// Resources
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

func (client *Client) DeleteResource(ctx context.Context, resourceID int32) (*resourcesClient.DeleteResourceReponse, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	deleteResourceRequest := &resourcesClient.DeleteResourceRequest{
		ResourceId: resourceID,
	}

	return client.resourcesService.DeleteResource(ctx, deleteResourceRequest)
}

func (client *Client) UpdateResource(ctx context.Context, resource *resourcesClient.Resource, scopes *ResourceScopeList) (*resourcesClient.ResourceData, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resourceData := &resourcesClient.ResourceData{
		Resource: resource,
		Scopes:   scopes.createResourceScopeRequestMapper().GetData(),
	}

	return client.resourcesService.UpdateResource(ctx, resourceData)
}

// Resource Scopes
func (client *Client) GetResourceScope(ctx context.Context, resourceID int32) (*resourcesClient.ListResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	getResourceScopeRequest := &resourcesClient.GetResourceScopeRequest{
		ResourceId: resourceID,
	}

	return client.resourcesService.GetResourceScope(ctx, getResourceScopeRequest)
}

// Clients
// func (client *Client) CreateClient(ctx context.Context, principal string, clientType string) (*clientsClient.Client, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	createClientRequest := &clientsClient.CreateClientRequest{
// 		Principal: principal,
// 		Type:      clientType,
// 	}

// 	return client.clientsService.CreateClient(ctx, createClientRequest)
// }

// func (client *Client) CreateRole(ctx context.Context, clientID int32, resourceScopeID int32) (*clientsClient.Role, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	clientRole := &clientsClient.CreateRoleRequest{
// 		ClientId:        clientID,
// 		ResourceScopeId: resourceScopeID,
// 	}

// 	return client.clientsService.CreateRole(ctx, clientRole)
// }

// func (client *Client) DeleteClient(ctx context.Context, clientID int32) (*clientsClient.DeleteClientReponse, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	deleteClientRequest := &clientsClient.DeleteClientRequest{
// 		ClientId: clientID,
// 	}

// 	return client.clientsService.DeleteClient(ctx, deleteClientRequest)
// }

// func (client *Client) DeleteRole(ctx context.Context, clientId int32, resourceScopeId int32) (*clientsClient.DeleteRoleReponse, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	deleteRoleRequest := &clientsClient.DeleteRoleRequest{
// 		ClientId:        clientId,
// 		ResourceScopeId: resourceScopeId,
// 	}

// 	return client.clientsService.DeleteRole(ctx, deleteRoleRequest)
// }

// func (client *Client) ListClients(ctx context.Context) (*clientsClient.ListClientResponse, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	return client.clientsService.ListClients(ctx, &emptypb.Empty{})
// }

// func (client *Client) ListClientScopes(ctx context.Context, clientID int32) (*resourcesClient.ListResourceScope, error) {
// 	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	listClientScopesRequest := &resourcesClient.ListClientScopesRequest{
// 		ClientId: clientID,
// 	}

// 	return client.resourcesService.ListClientScopes(ctx, listClientScopesRequest)
// }
