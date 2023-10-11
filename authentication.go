package indicoserviceauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	authClient "github.com/indicoinnovation/indico_service_auth/client/auth"
	integrationClient "github.com/indicoinnovation/indico_service_auth/client/integrations"
	mfaClient "github.com/indicoinnovation/indico_service_auth/client/mfa"

	"github.com/go-jose/go-jose"
	"github.com/google/uuid"

	"github.com/indicoinnovation/indico_service_auth/pkg/constants"
	"github.com/indicoinnovation/indico_service_auth/pkg/helpers"
	"github.com/indicoinnovation/indico_service_auth/pkg/iam"
)

type Client struct {
	authService        authClient.AuthServiceClient
	mfaService         mfaClient.MFAServiceClient
	integrationService integrationClient.IntegrationsServiceClient
}

func generateToken(scope string) (string, error) {
	publicKey, err := helpers.ParsePublicFromPrivateKey(iam.Credentials.PrivateKey)
	if err != nil {
		return "", err
	}

	encrypter, err := jose.NewEncrypter(jose.A256CBC_HS512, jose.Recipient{Algorithm: jose.RSA_OAEP, KeyID: iam.Credentials.PrivateKeyID, Key: publicKey}, nil)
	if err != nil {
		return "", err
	}

	claims, err := json.Marshal(map[string]interface{}{
		"type":      iam.Credentials.Type,
		"principal": iam.Credentials.Principal,
		"scope":     scope,
		"svc":       strings.Split(iam.Credentials.Principal, "@")[0],
		"aud":       constants.Audience,
		"jti":       uuid.New().String(),
		"exp":       time.Now().Unix() + constants.AccessTokenValidity,
		"iat":       time.Now().Unix(),
		"iss":       iam.Credentials.Principal,
		"sub":       iam.Credentials.Principal,
	})
	if err != nil {
		return "", err
	}

	plaintext := []byte(claims)
	jwe, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return "", err
	}

	serialized, err := jwe.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}

func authenticate(context context.Context, authservice authClient.AuthServiceClient, accessToken string) (*authClient.AuthToken, error) {
	request := &authClient.AuthToken{
		Jwt: accessToken,
	}

	return authservice.Authenticate(context, request)
}

func authorize(context context.Context, client *Client, scope string) error {
	resource := strings.Split(scope, ".")[0]
	// desirableScope := strings.Split(scope, ".")[1]
	adminScope := fmt.Sprintf("%s.admin", resource)
	scopes := []string{scope, adminScope, constants.IAMAdmin}

	for _, currentScope := range scopes {
		token, err := generateToken(currentScope)
		if err != nil {
			return fmt.Errorf("error to generate jwt token: %w", err)
		}

		authenticated, err := authenticate(context, client.authService, token)
		if err != nil {
			return fmt.Errorf("error to authenticate: %s", err)
		}

		if authenticated.Jwt != "" {
			return nil
		}
	}

	return fmt.Errorf("%w", errors.New("service account not allowed to make this request. check permissions and try again"))
}

func NewClient() (*Client, context.Context, error) {
	ctx, _ := helpers.InitContext()
	conn := iam.Connect()
	client := &Client{
		authService:        authClient.NewAuthServiceClient(conn),
		mfaService:         mfaClient.NewMFAServiceClient(conn),
		integrationService: integrationClient.NewIntegrationsServiceClient(conn),
	}

	err := authorize(ctx, client, constants.AuthConnect)

	return client, ctx, err
}
