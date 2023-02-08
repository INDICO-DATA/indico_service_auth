package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	indicoserviceauth "github.com/INDICO-INNOVATION/indico_service_auth"
	"github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	resource "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	resourceclient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceaccounts "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"
)

func main() {
	client, ctx, err := indicoserviceauth.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}

	testBackoffice(client, ctx)

	// generateAndValidate(client, ctx)
	// validateThird(client, ctx, "643863")
}

// func generateAndValidate(client *indicoserviceauth.Client, ctx context.Context) {
// 	response, err := client.GenerateOTP(ctx)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Generate OTP Response:")
// 	fmt.Printf("%+v\n\n", response)

// 	responsev, err := client.ValidateOTP(ctx, response.Token, true)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Validate OTP Response:")
// 	fmt.Printf("%+v\n\n", responsev)
// }

// func validateThird(client *indicoserviceauth.Client, ctx context.Context, token string) {
// 	responsev, err := client.ValidateOTP(ctx, token, false)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Validate Third Party OTP Response:")
// 	fmt.Printf("%+v\n\n", responsev)

// 	ctx.Done()
// }

func testBackoffice(client *indicoserviceauth.Client, ctx context.Context) {
	testServiceAccounts(client, ctx)
	testClients(client, ctx)
	testResources(client, ctx)
}

func testServiceAccounts(client *indicoserviceauth.Client, ctx context.Context) {
	var serviceAccounts []*serviceaccounts.ServiceAccount
	var serviceAccountsToDelete []string

	for _, i := range []string{"1", "2", "3"} {
		serviceAccount, err := client.CreateServiceAccount(ctx, "Teste"+i, "teste"+i, "sa creation lib test")
		if err != nil {
			log.Fatalf(err.Error())
		}
		serviceAccounts = append(serviceAccounts, serviceAccount)
		serviceAccountsToDelete = append(serviceAccountsToDelete, strconv.Itoa(int(serviceAccount.ServiceAccountId)))
	}

	fmt.Println("Created SAs:")
	fmt.Println(serviceAccounts)
	fmt.Println("")

	serviceAccounts, err := client.DeleteServiceAccount(ctx, serviceAccountsToDelete)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Deleted SA:")
	fmt.Printf("%+v", serviceAccounts)
	fmt.Println("")

	serviceAccounts, err = client.ListServiceAccounts(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("SA Listing:")
	fmt.Printf("%+v \n", serviceAccounts)
	fmt.Println("")

	credentials, err := client.GenerateCredentials(ctx, strconv.Itoa(int(serviceAccounts[0].ServiceAccountId)))
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Credentials Generated:")
	fmt.Printf("%+v \n", credentials)
	fmt.Println("")
}

func testClients(client *indicoserviceauth.Client, ctx context.Context) {
	var clients []*clients.Client
	var clientsToDelete []int64

	for _, i := range []string{"indico_user@indico.pt", "indico_user@gmail.pt", "indico_user@outlook.pt"} {
		client, err := client.CreateClient(ctx, i, "user")
		if err != nil {
			log.Fatalf(err.Error())
		}

		clients = append(clients, client)
		clientsToDelete = append(clientsToDelete, client.Id)
	}

	fmt.Println("Created clients:")
	fmt.Println(clients)
	fmt.Println("")

	updatedClient, err := client.UpdateClient(ctx, clients[0].Id, "indico_user_updated@indico.pt")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Client updated:")
	fmt.Printf("%+v \n", updatedClient)
	fmt.Println("")

	roles := []*clientsClient.Role{
		{ScopeId: 1},
		{ScopeId: 2},
	}

	roles, err = client.CreateRole(ctx, clients[0].Id, roles)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Roles created:")
	fmt.Printf("%+v \n", roles)
	fmt.Println("")

	retrievedScopes, err := client.RetrieveClient(ctx, clients[0].Id)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Client retrieved scopes:")
	fmt.Printf("%+v \n", retrievedScopes)
	fmt.Println("")

	roles[0].ScopeId = 3
	updatedRoles, err := client.UpdateRoles(ctx, roles)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Roles updated:")
	fmt.Printf("%+v \n", updatedRoles)
	fmt.Println("")

	clients, err = client.ListClients(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Clients Listing:")
	fmt.Printf("%+v \n", clients)
	fmt.Println("")

	clients, err = client.DeleteClient(ctx, clientsToDelete)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Deleted clients:")
	fmt.Printf("%+v", clients)
	fmt.Println("")
}

func testResources(client *indicoserviceauth.Client, ctx context.Context) {
	// Create Resources
	fmt.Println("Created Resources")

	var resources []*resourceclient.Resource
	var resourcesToDelete []string

	for _, i := range []string{"recurso_teste1", "recurso_teste2", "recurso_teste3"} {
		resource, err := client.CreateResource(ctx, i, "description")
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		resources = append(resources, resource)
		resourcesToDelete = append(resourcesToDelete, strconv.Itoa(int(resource.ResourceId)))
	}

	fmt.Printf("%+v\n\n", resources)

	// List Resources
	fmt.Println("Resources Listing")

	resourceList, err := client.ListResources(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n\n", resourceList)

	// Update Resource
	fmt.Println("Updated resource")

	resourceID, _ := strconv.Atoi(resourcesToDelete[0])

	resourcesToUpdate := &resource.Resource{
		ResourceId:  int64(resourceID),
		Name:        "resourcesname",
		Description: "resources[0].Description",
	}

	resource, err := client.UpdateResource(ctx, resourcesToUpdate)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n", resource)

	// Create Scope
	fmt.Println("Create scopes")

	scopes := []*resourceclient.ResourceScope{
		{ResourceId: int64(resourceID), Name: "editor", Label: "Editor", Description: "Editar"},
	}
	scope, err := client.CreateResourceScope(ctx, scopes)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n", scope)

	// Update Scopes
	fmt.Println("Updated scopes")

	scopesToUpdate := []*resourceclient.ResourceScope{
		{ResourceId: resourceList[0].ResourceId, Name: "update", Label: "Teste update", Description: "update_desc"},
		{ResourceId: resourceList[0].ResourceId, Name: "update_two", Label: "Teste update Two", Description: "update_desc_two"},
	}

	scopes, err = client.UpdateResourceScope(ctx, scopesToUpdate)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n", scopes)

	// Delete Resources
	fmt.Println("Deleted resources")

	deletedResources, err := client.DeleteResource(ctx, resourcesToDelete)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n\n", deletedResources)
}
