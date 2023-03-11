package indicoserviceauth

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	resourcesClient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceAccountClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"
	serviceAccountKeyClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account_keys"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
		if err := stream.Send(&serviceAccountClient.ServiceAccountRequest{ServiceAccountId: accountID}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming delete service accounts: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return deletedServicesAccounts, nil
}

// Service Account Keys
func (client *Client) CreateServiceAccountKey(ctx context.Context, serviceAccountID int64) (*serviceAccountKeyClient.ServiceAccountKey, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.serviceAccountKeysService.Create(ctx, &serviceAccountKeyClient.CredentialsRequest{ServiceAccountId: serviceAccountID})
}

func (client *Client) ListServiceAccountKeys(ctx context.Context, serviceAccountID int64) ([]*serviceAccountKeyClient.ServiceAccountKey, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keys := []*serviceAccountKeyClient.ServiceAccountKey{}

	stream, err := client.serviceAccountKeysService.List(ctx, &serviceAccountKeyClient.CredentialsRequest{ServiceAccountId: serviceAccountID})
	if err != nil {
		return nil, fmt.Errorf("error to stream list service account keys due to: %w", err)
	}

	for {
		serviceAccount, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming list service account keys: %w", err)
		}

		keys = append(keys, serviceAccount)
	}

	return keys, nil
}

func (client *Client) DeleteServiceAccountKey(ctx context.Context, keyIDs []string) ([]*serviceAccountKeyClient.ServiceAccountKey, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.serviceAccountKeysService.Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming delete service account keys: %w", err)
	}

	deletedServicesAccounts := []*serviceAccountKeyClient.ServiceAccountKey{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming delete service account keys response: %s", err.Error())
			}

			deletedServicesAccounts = append(deletedServicesAccounts, in)
		}
	}()

	for _, key := range keyIDs {
		if err := stream.Send(&serviceAccountKeyClient.ServiceAccountKeyRequest{KeyId: key}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming delete service account keys: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return deletedServicesAccounts, nil
}

func (client *Client) RetrieveServiceAccountKey(ctx context.Context, keyID string) (*serviceAccountKeyClient.ServiceAccountKey, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.serviceAccountKeysService.Retrieve(ctx, &serviceAccountKeyClient.ServiceAccountKeyRequest{KeyId: keyID})
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

	return client.resourcesService.Create(ctx, createResourceRequest)
}

func (client *Client) ListResources(ctx context.Context) ([]*resourcesClient.Resource, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resources := []*resourcesClient.Resource{}

	stream, err := client.resourcesService.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("error to stream list service account due to: %w", err)
	}

	for {
		resource, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming list resources: %w", err)
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (client *Client) DeleteResource(ctx context.Context, resourceIDs []string) ([]*resourcesClient.Resource, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.resourcesService.Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming delete resources: %w", err)
	}

	deletedResources := []*resourcesClient.Resource{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming delete resources response: %s", err.Error())
			}

			deletedResources = append(deletedResources, in)
		}
	}()

	for _, resourceID := range resourceIDs {
		if err := stream.Send(&resourcesClient.QueryResourceRequest{ResourceId: resourceID}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming delete resources: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return deletedResources, nil
}

func (client *Client) UpdateResource(ctx context.Context, resource *resourcesClient.Resource) (*resourcesClient.Resource, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.resourcesService.Update(ctx, resource)
}

// Resource Scopes
func (client *Client) CreateResourceScope(ctx context.Context, scopes []*resourcesClient.ResourceScope) ([]*resourcesClient.ResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.resourcesService.CreateScope(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming delete resources: %w", err)
	}

	var resourceScopes []*resourcesClient.ResourceScope

	waitc := make(chan struct{})
	go func() {
		for {
			resourceScope, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming create resource scopes response: %s", err.Error())
			}

			resourceScopes = append(resourceScopes, resourceScope)
		}
	}()

	for _, scope := range scopes {
		err := stream.Send(&resourcesClient.ResourceScope{
			Name:        scope.Name,
			Label:       scope.Label,
			Description: scope.Description,
			ResourceId:  int64(scope.ResourceId),
		})

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming create resource scopes: %w", err)
		}
	}

	stream.CloseSend()
	<-waitc

	return resourceScopes, nil
}

func (client *Client) ListResourceScope(ctx context.Context, resourceID string) ([]*resourcesClient.ResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	getResourceScopeRequest := &resourcesClient.QueryResourceRequest{
		ResourceId: resourceID,
	}

	stream, err := client.resourcesService.ListScope(ctx, getResourceScopeRequest)
	if err != nil {
		return nil, fmt.Errorf("error to stream list resource scpopes due to: %w", err)
	}

	scopes := []*resourcesClient.ResourceScope{}

	for {
		scope, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming list resources: %w", err)
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}

func (client *Client) UpdateResourceScope(ctx context.Context, scopes []*resourcesClient.ResourceScope) ([]*resourcesClient.ResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.resourcesService.UpdateScope(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming update resource scopes: %w", err)
	}

	updatedResources := []*resourcesClient.ResourceScope{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming update resource scopes response: %s", err.Error())
			}

			updatedResources = append(updatedResources, in)
		}
	}()

	for _, scope := range scopes {
		if err := stream.Send(&resourcesClient.ResourceScope{
			ResourceScopesId: scope.ResourceScopesId,
			Label:            scope.Label,
			Name:             scope.Name,
			Description:      scope.Description,
			ResourceId:       scope.ResourceId,
		}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming update resource scopes: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return updatedResources, nil
}

// Clients
func (client *Client) CreateClient(ctx context.Context, principal string, clientType string) (*clientsClient.Client, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createClientRequest := &clientsClient.ClientRequest{
		Principal: principal,
		Type:      clientType,
	}

	return client.clientsService.Create(ctx, createClientRequest)
}

func (client *Client) CreateRole(ctx context.Context, clientID int64, roles []*clientsClient.Role) ([]*clientsClient.Role, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.clientsService.CreateRole(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming creating roles: %w", err)
	}

	createdRoles := []*clientsClient.Role{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming delete creating roles: %s", err.Error())
			}

			createdRoles = append(roles, in)
		}
	}()

	for _, role := range roles {
		if err := stream.Send(&clientsClient.CreateRoleRequest{ClientId: clientID, ScopeId: role.ScopeId}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming creating roles: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return createdRoles, nil
}

func (client *Client) DeleteClient(ctx context.Context, clientIDs []int64) ([]*clientsClient.Client, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.clientsService.Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming delete clients: %w", err)
	}

	deletedClients := []*clientsClient.Client{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming delete clients response: %s", err.Error())
			}

			deletedClients = append(deletedClients, in)
		}
	}()

	for _, clientID := range clientIDs {
		if err := stream.Send(&clientsClient.DeleteClientRequest{ClientId: clientID}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming delete clients: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return deletedClients, nil
}

func (client *Client) ListClients(ctx context.Context) ([]*clientsClient.Client, error) {
	if err := authorize(ctx, client, "iam_server.use"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	clients := []*clientsClient.Client{}

	stream, err := client.clientsService.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("error to stream list clients due to: %w", err)
	}

	for {
		client, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming list clients: %w", err)
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func (client *Client) UpdateClient(ctx context.Context, id int64, principal string) (*clientsClient.Client, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	createClientRequest := &clientsClient.Client{
		Id:        id,
		Principal: principal,
	}

	return client.clientsService.Update(ctx, createClientRequest)
}

func (client *Client) UpdateRoles(ctx context.Context, roles []*clientsClient.Role) ([]*clientsClient.Role, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream, err := client.clientsService.UpdateRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("the following error occured while streaming updating roles: %w", err)
	}

	updatedRoles := []*clientsClient.Role{}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("the following error occured while streaming updating roles: %s", err.Error())
			}

			updatedRoles = append(roles, in)
		}
	}()

	for _, role := range roles {
		if err := stream.Send(&clientsClient.Role{ClientId: role.ClientId, ScopeId: role.ScopeId, Id: role.Id}); err != nil {
			return nil, fmt.Errorf("the following error occured while streaming updating roles: %w", err)
		}
	}
	stream.CloseSend()
	<-waitc

	return updatedRoles, nil
}

func (client *Client) RetrieveClient(ctx context.Context, clientID int64) ([]*clientsClient.RetrieveResourceScope, error) {
	if err := authorize(ctx, client, "iam_backoffice.admin"); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	retrieveRequest := clientsClient.DeleteClientRequest{ClientId: clientID}

	scopes := []*clients.RetrieveResourceScope{}

	stream, err := client.clientsService.Retrieve(ctx, &retrieveRequest)
	if err != nil {
		return nil, fmt.Errorf("error to stream retrieve client due to: %w", err)
	}

	for {
		scope, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("the following error occured while streaming retrieve client: %w", err)
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}
