package main

import (
	"context"
	"fmt"
	"log"
	"os"

	indicoserviceauth "github.com/INDICO-INNOVATION/indico_service_auth/v2"
)

func main() {
	//testAPIClients()

	testBackoffice()

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

func testAPIClients() {
	log.Printf("Testing API Clients credentials to IAM...\n\n")
	// testIamBackofficeAPIClient()
	testGaiaAPIClient()
}

func testIamBackofficeAPIClient() {
	os.Setenv("INNOVATION_CREDENTIALS", "backoffice.json")

	_, _, err := indicoserviceauth.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("Backoffice API Client successfully connected to IAM through Indico Service Auth")
}

func testGaiaAPIClient() {
	os.Setenv("INNOVATION_CREDENTIALS", "gaia.json")

	client, ctx, err := indicoserviceauth.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("G.A.I.A API Client successfully connected to IAM through Indico Service Auth")

	headers := map[string]interface{}{
		"type": "JWT",
		"alg":  "RS512",
	}

	claims := map[string]interface{}{
		"aud": "https://iam.services.indicoinnovation.pt",
		"jti": "jti",
		"exp": 2677868842,
		"iat": 2677868842,
		"iss": "iss",
		"sub": "iss",
	}

	private := "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCdPgd3d+/v4Hen\n95n2N+RP6K8WGL1W9isxmpS0AcXThneiLXv7Idy1xyjQdPhog4AS7M//eobkD6OK\n4udjioR7cBzZqdDV3pI1/ArkrLwxu+kldwjx1vFhO2O8vcKw9+nru1CgZve65Cv+\n12/1KleSyeqeicYGQh0z6tliqrlg7ySXJOuJMwvkcYygJ/ehKQwzWKw3sfUOuhmu\nUH2b8L64OAeOSoas5yi7nY53BRRBkYJSGhen/gIULBXfvdWTQITILI60GZq0DHsB\nQ+Q6FPLvLlzbhcqVcFcF6EHjPM4FBgFm9llshtOcwHDWDGDocpGG5PzyMZQQFkiZ\naENiokQVAgMBAAECggEABRGjf/sX0Fh2GFi9AoOZZ+diWDNX9zOxqhmyQtNDzfvU\nY3VVPSqW84KGEveSV8bVcXFnix7SC7a4A5l+SBExECGfFyKISZIGje2gFEW/6K0x\naN8begr/IQP0b83WNt+RujhW8woetkxLYXMPIuId4ezXFAvw8ZwJVkvIfOEatMCJ\nKobd3x3BsEnINmM6ARBZ6nG6x0jxAizE7o0ZiFk1702Dzx+PCEHdnq832j1nlCe1\n18Q7Fw+cBP8mGwemuWJmGkhi9t6oNBLytuWfp0lEo2TMg5JGcF5zAaoAOZCHEiAH\n3g6T17c9Vhx3W0x7wTXVgTxyzIc3lfx6M2BK0jxISQKBgQDXIc77uh9bWJ+jXpk5\nJn370eRQ3NVQHLdPjHaXnVnH7XlLGPTxZEY/+gQnw3lKfkS/bbh0VVfPN7RkJAuA\nLzfhyVKbuaIzkH2vfC39HaoHOwFMmiMz6u6/AbMEyaOwXgrwNvUiGgBYQu7lOUm3\nPpyrFAzJ0SKgQtK1xM/royFN7wKBgQC7HPSgUia80deBWWAqc0Yg7O2f+lxTdkdm\nxW6b5ydJxjUHWOUIvosIMd3CSy8nChIdBzCUdsFKnl2hJq6luYb2R9HM8Goa8/+2\n27lruf5bq0cBghO+VUD2cALssXV+l9JahwjSY66A14pHYSmIawlaUp4xjoRQzF4g\n3S3Y762SOwKBgQC+GHvDZ6WJm/gD7fVK0L6GI+OxqFE9ZwJ74kxwkHPi++4v86+N\nFeQVaCD/2fR1ZLnuxU82TKs4yao2yrr228JvWuwJ+rBfDBBjXkFs5id2GCd/vrvY\nhoHMef3r/xkLIS5dMNLO/VbMpuhd2I9QwpzBdwgN1STZo/J4954aQ52MNQKBgCSl\nGBep3yJNqnoPqfDGYK6kGCmjm59Q48zxz6bCz4P4SHxm4Xwj7RCuy3J4zQRFSaTb\nupzo4RQfDdE4xS15tz/WcvVDBFTBAoyqXZcGdxU25xIZiopVIgLhjwESBCPF0hGe\nrpCdqxhanh86nSq1Y8CHu6mS8sDsNWcUSh1ZW1HvAoGAApVNy6f8JLebUC1lPTH8\nDRV6bmIsi606aKfq3Xzr7rCJvf9Qv0u0FtsAdmnE4mS7EEhnHARB9nLtOjHZ+xPj\noD8F8A6CL+fYfuLwuIOmkH3wijtRslcBcL+Y5AH+ydO0gJ+rNh7PgE16Pgh/I6GL\n+86HwCUd8uR75K7bDI6B/LA=\n-----END PRIVATE KEY-----\n"

	token, err := client.GenerateJWT(ctx, headers, claims, private)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println(token.Jwt)

	_, err = client.ValidateJWT(ctx, token.Jwt, private)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("token successfully validated")

}

func testBackoffice() {
	os.Setenv("INNOVATION_CREDENTIALS", "backoffice.json")

	client, ctx, err := indicoserviceauth.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// testServiceAccounts(client, ctx)
	testServiceAccountKeys(client, ctx)
	// testClients(client, ctx)
	// testResources(client, ctx)
}

// func testServiceAccounts(client *indicoserviceauth.Client, ctx context.Context) {
// 	var serviceAccounts []*serviceaccounts.ServiceAccount
// 	var serviceAccountsToDelete []string

// 	for _, i := range []string{"1", "2", "3"} {
// 		serviceAccount, err := client.CreateServiceAccount(ctx, "Teste"+i, "teste"+i, "sa creation lib test")
// 		if err != nil {
// 			log.Fatalf(err.Error())
// 		}
// 		serviceAccounts = append(serviceAccounts, serviceAccount)
// 		serviceAccountsToDelete = append(serviceAccountsToDelete, strconv.Itoa(int(serviceAccount.ServiceAccountId)))
// 	}

// 	fmt.Println("Created SAs:")
// 	fmt.Println(serviceAccounts)
// 	fmt.Println("")

// 	serviceAccounts, err := client.DeleteServiceAccount(ctx, serviceAccountsToDelete)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Deleted SA:")
// 	fmt.Printf("%+v", serviceAccounts)
// 	fmt.Println("")

// 	serviceAccounts, err = client.ListServiceAccounts(ctx)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("SA Listing:")
// 	fmt.Printf("%+v \n", serviceAccounts)
// 	fmt.Println("")

// }

func testServiceAccountKeys(client *indicoserviceauth.Client, ctx context.Context) {
	serviceAccounts, err := client.ListServiceAccounts(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(serviceAccounts) == 0 {
		log.Fatal("no service accounts found to test service account keys")
	}

	key, err := client.CreateServiceAccountKey(ctx, serviceAccounts[0].ServiceAccountId)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Service Account Key Generated:")
	fmt.Printf("%+v \n", key)
	fmt.Println("")

	keys, err := client.ListServiceAccountKeys(ctx, serviceAccounts[0].ServiceAccountId)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Service Account Key Listed:")
	fmt.Printf("%+v \n", keys)
	fmt.Println("")

	key, err = client.RetrieveServiceAccountKey(ctx, key.KeyId)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Service Account Key Retrieved:")
	fmt.Printf("%+v \n", key)
	fmt.Println("")

	keys, err = client.DeleteServiceAccountKey(ctx, []string{key.KeyId})
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Service Account Key Deleted:")
	fmt.Printf("%+v \n", keys)
	fmt.Println("")
}

// func testClients(client *indicoserviceauth.Client, ctx context.Context) {
// 	var clients []*clients.Client
// 	var clientsToDelete []int64

// 	for _, i := range []string{"indico_user@indico.pt", "indico_user@gmail.pt", "indico_user@outlook.pt"} {
// 		client, err := client.CreateClient(ctx, i, "user")
// 		if err != nil {
// 			log.Fatalf(err.Error())
// 		}

// 		clients = append(clients, client)
// 		clientsToDelete = append(clientsToDelete, client.Id)
// 	}

// 	fmt.Println("Created clients:")
// 	fmt.Println(clients)
// 	fmt.Println("")

// 	updatedClient, err := client.UpdateClient(ctx, clients[0].Id, "indico_user_updated@indico.pt")
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Client updated:")
// 	fmt.Printf("%+v \n", updatedClient)
// 	fmt.Println("")

// 	roles := []*clientsClient.Role{
// 		{ScopeId: 1},
// 		{ScopeId: 2},
// 	}

// 	roles, err = client.CreateRole(ctx, clients[0].Id, roles)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Roles created:")
// 	fmt.Printf("%+v \n", roles)
// 	fmt.Println("")

// 	retrievedScopes, err := client.RetrieveClient(ctx, clients[0].Id)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Client retrieved scopes:")
// 	fmt.Printf("%+v \n", retrievedScopes)
// 	fmt.Println("")

// 	roles[0].ScopeId = 3
// 	updatedRoles, err := client.UpdateRoles(ctx, roles)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Roles updated:")
// 	fmt.Printf("%+v \n", updatedRoles)
// 	fmt.Println("")

// 	clients, err = client.ListClients(ctx)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Clients Listing:")
// 	fmt.Printf("%+v \n", clients)
// 	fmt.Println("")

// 	clients, err = client.DeleteClient(ctx, clientsToDelete)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Println("Deleted clients:")
// 	fmt.Printf("%+v", clients)
// 	fmt.Println("")
// }

// func testResources(client *indicoserviceauth.Client, ctx context.Context) {
// 	// Create Resources
// 	fmt.Println("Created Resources")

// 	var resources []*resourceclient.Resource
// 	var resourcesToDelete []string

// 	for _, i := range []string{"recurso_teste1", "recurso_teste2", "recurso_teste3"} {
// 		resource, err := client.CreateResource(ctx, i, "description")
// 		if err != nil {
// 			log.Fatalf("%s", err.Error())
// 		}

// 		resources = append(resources, resource)
// 		resourcesToDelete = append(resourcesToDelete, strconv.Itoa(int(resource.ResourceId)))
// 	}

// 	fmt.Printf("%+v\n\n", resources)

// 	// List Resources
// 	fmt.Println("Resources Listing")

// 	resourceList, err := client.ListResources(ctx)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Printf("%+v\n\n", resourceList)

// 	// Update Resource
// 	fmt.Println("Updated resource")

// 	resourceID, _ := strconv.Atoi(resourcesToDelete[0])

// 	resourcesToUpdate := &resource.Resource{
// 		ResourceId:  int64(resourceID),
// 		Name:        "resourcesname",
// 		Description: "resources[0].Description",
// 	}

// 	resource, err := client.UpdateResource(ctx, resourcesToUpdate)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Printf("%+v\n", resource)

// 	// Create Scope
// 	fmt.Println("Create scopes")

// 	scopes := []*resourceclient.ResourceScope{
// 		{ResourceId: int64(resourceID), Name: "editor", Label: "Editor", Description: "Editar"},
// 	}
// 	scope, err := client.CreateResourceScope(ctx, scopes)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Printf("%+v\n", scope)

// 	// Update Scopes
// 	fmt.Println("Updated scopes")

// 	scopesToUpdate := []*resourceclient.ResourceScope{
// 		{ResourceId: resourceList[0].ResourceId, Name: "update", Label: "Teste update", Description: "update_desc"},
// 		{ResourceId: resourceList[0].ResourceId, Name: "update_two", Label: "Teste update Two", Description: "update_desc_two"},
// 	}

// 	scopes, err = client.UpdateResourceScope(ctx, scopesToUpdate)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Printf("%+v\n", scopes)

// 	// Delete Resources
// 	fmt.Println("Deleted resources")

// 	deletedResources, err := client.DeleteResource(ctx, resourcesToDelete)
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	fmt.Printf("%+v\n\n", deletedResources)
// }
