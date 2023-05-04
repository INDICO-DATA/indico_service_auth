package indicoserviceauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	authClient "github.com/INDICO-INNOVATION/indico_service_auth/client/auth"
	clientsClient "github.com/INDICO-INNOVATION/indico_service_auth/client/clients"
	integrationClient "github.com/INDICO-INNOVATION/indico_service_auth/client/integrations"
	mfaClient "github.com/INDICO-INNOVATION/indico_service_auth/client/mfa"
	resourcesClient "github.com/INDICO-INNOVATION/indico_service_auth/client/resources"
	serviceAccountClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account"
	serviceAccountKeysClient "github.com/INDICO-INNOVATION/indico_service_auth/client/service_account_keys"

	"github.com/go-jose/go-jose"
	"github.com/google/uuid"

	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/constants"
	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/helpers"
	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/iam"
)

type Client struct {
	authService               authClient.AuthServiceClient
	mfaService                mfaClient.MFAServiceClient
	serviceAccountService     serviceAccountClient.ServiceAccountServiceClient
	serviceAccountKeysService serviceAccountKeysClient.ServiceAccountKeysServiceClient
	resourcesService          resourcesClient.ResourceServiceClient
	clientsService            clientsClient.ClientServiceClient
	integrationService        integrationClient.IntegrationsServiceClient
}

func generateToken(ctx context.Context, authservice authClient.AuthServiceClient, scope string) (*authClient.AuthToken, error) {
	publicKey, err := helpers.ParsePublicFromPrivateKey(iam.Credentials.PrivateKey)
	if err != nil {
		return &authClient.AuthToken{Jwt: ""}, err
	}

	encrypter, err := jose.NewEncrypter(jose.A256CBC_HS512, jose.Recipient{Algorithm: jose.RSA_OAEP, KeyID: iam.Credentials.PrivateKeyID, Key: publicKey}, nil)
	if err != nil {
		return &authClient.AuthToken{Jwt: ""}, err
	}

	clains, err := json.Marshal(map[string]interface{}{
		"type":      iam.Credentials.Type,
		"principal": iam.Credentials.Principal,
		"scope":     scope,
		"svc":       strings.Split(iam.Credentials.Principal, "@")[0],
		"aud":       constants.AUDIENCE,
		"jti":       uuid.New().String(),
		"exp":       time.Now().Unix() + constants.ACCESS_TOKEN_VALIDITY,
		"iat":       time.Now().Unix(),
		"iss":       iam.Credentials.Principal,
		"sub":       iam.Credentials.Principal,
	})
	if err != nil {
		return &authClient.AuthToken{Jwt: ""}, err
	}

	plaintext := []byte(clains)
	jwe, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return &authClient.AuthToken{Jwt: ""}, err
	}

	serialized, err := jwe.CompactSerialize()
	if err != nil {
		return &authClient.AuthToken{Jwt: ""}, err
	}

	return &authClient.AuthToken{Jwt: serialized}, nil
}

func authenticate(context context.Context, authservice authClient.AuthServiceClient, accessToken string) (*authClient.AuthToken, error) {
	request := &authClient.AuthToken{
		Jwt: accessToken,
	}

	return authservice.Authenticate(context, request)
}

func authorize(context context.Context, client *Client, scope string) error {
	response, err := client.clientsService.IsOwner(context, &clientsClient.IsOwnerRequest{Principal: iam.Credentials.Principal})
	if err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("error to verify if user is owner: %s", err.Error()))
	}

	if response.IsOwner {
		return nil
	}

	token, err := generateToken(context, client.authService, scope)
	if err != nil {
		return fmt.Errorf("error to generate jwt token: %w", err)
	}

	authenticated, err := authenticate(context, client.authService, token.Jwt)
	if err != nil {
		return fmt.Errorf("error to authenticate: %s", err)
	}
	if authenticated.Jwt == "" {
		return fmt.Errorf("%w", errors.New("you are not allowed to make this request"))
	}

	return nil
}

func NewClient() (*Client, context.Context, error) {
	ctx, _ := helpers.InitContext()
	conn := iam.Connect()
	client := &Client{
		authService:               authClient.NewAuthServiceClient(conn),
		mfaService:                mfaClient.NewMFAServiceClient(conn),
		serviceAccountService:     serviceAccountClient.NewServiceAccountServiceClient(conn),
		resourcesService:          resourcesClient.NewResourceServiceClient(conn),
		clientsService:            clientsClient.NewClientServiceClient(conn),
		integrationService:        integrationClient.NewIntegrationsServiceClient(conn),
		serviceAccountKeysService: serviceAccountKeysClient.NewServiceAccountKeysServiceClient(conn),
	}

	err := authorize(ctx, client, "auth.connect")

	return client, ctx, err
}
