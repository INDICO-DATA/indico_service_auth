package indicoserviceauth

import (
	"context"
	"encoding/json"
	"fmt"

	authClient "github.com/indicoinnovation/indico_service_auth/client/auth"
	integrationClient "github.com/indicoinnovation/indico_service_auth/client/integrations"
	"github.com/indicoinnovation/indico_service_auth/pkg/constants"
)

func (client *Client) GenerateJWT(ctx context.Context, headers map[string]interface{}, claims map[string]interface{}, private string) (*authClient.AuthToken, error) {
	if err := authorize(ctx, client, constants.TokenGenerate); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	headerBytes, err := json.Marshal(headers)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.integrationService.GenerateJWT(ctx, &integrationClient.GenerateJWTRequest{
		Header:     string(headerBytes),
		Claims:     string(claimsBytes),
		PrivateKey: private,
	})
}

func (client *Client) ValidateJWT(ctx context.Context, token string, privateKey string) (*authClient.AuthToken, error) {
	if err := authorize(ctx, client, constants.TokenValidate); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return client.integrationService.ValidateJWT(ctx, &integrationClient.ValidateJWTRequest{Jwt: token, PrivateKey: privateKey})
}
