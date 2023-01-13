package indicoserviceauth

import (
	"context"
	"errors"
	"fmt"

	authClient "github.com/INDICO-INNOVATION/indico_service_auth/client/auth"
	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	mfaClient "github.com/INDICO-INNOVATION/indico_service_auth/client/mfa"
	resourcesClient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceAccountClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"

	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/helpers"
	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/iam"
)

type Client struct {
	authService           authClient.AuthServiceClient
	mfaService            mfaClient.MFAServiceClient
	serviceAccountService serviceAccountClient.ServiceAccountServiceClient
	resourcesService      resourcesClient.ResourceServiceClient
	clientsService        clientsClient.ClientServiceClient
}

func generateToken(context context.Context, authservice authClient.AuthServiceClient, scope string) (*authClient.GenerateTokenResponse, error) {
	request := &authClient.GenerateTokenRequest{
		ClientId:     iam.Credentials.ClientID,
		ClientSecret: iam.Credentials.ClientSecret,
		Scope:        scope,
		PrivateKey:   iam.Credentials.PrivateKey,
		Type:         iam.Credentials.Type,
	}

	return authservice.GenerateToken(context, request)
}

func authenticate(context context.Context, authservice authClient.AuthServiceClient, accessToken string) (*authClient.AuthResponse, error) {
	request := &authClient.AuthRequest{
		AccessToken: accessToken,
		PrivateKey:  iam.Credentials.PrivateKey,
	}

	return authservice.Authenticate(context, request)
}

func authorize(context context.Context, client *Client, scope string) error {
	token, err := generateToken(context, client.authService, scope)
	if err != nil {
		return fmt.Errorf("error to generate jwt token: %w", err)
	}

	authenticated, err := authenticate(context, client.authService, token.AccessToken)
	if err != nil {
		return fmt.Errorf("error to authenticate: %s", err)
	}
	if !authenticated.Authenticated {
		return fmt.Errorf("%w", errors.New("you are not allowed to make this request"))
	}
	return nil
}

func NewClient() (*Client, context.Context, error) {
	ctx, _ := helpers.InitContext()
	conn := iam.Connect()

	client := &Client{
		authService:           authClient.NewAuthServiceClient(conn),
		mfaService:            mfaClient.NewMFAServiceClient(conn),
		serviceAccountService: serviceAccountClient.NewServiceAccountServiceClient(conn),
		resourcesService:      resourcesClient.NewResourceServiceClient(conn),
		clientsService:        clientsClient.NewClientServiceClient(conn),
	}

	err := authorize(ctx, client, "auth.connect")

	return client, ctx, err
}
