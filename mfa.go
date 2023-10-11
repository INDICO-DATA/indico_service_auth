package indicoserviceauth

import (
	"context"
	"fmt"

	mfaClient "github.com/indicoinnovation/indico_service_auth/client/mfa"
	"github.com/indicoinnovation/indico_service_auth/pkg/constants"
	"github.com/indicoinnovation/indico_service_auth/pkg/iam"
)

func (client *Client) GenerateOTP(ctx context.Context, clientSecret string) (*mfaClient.GenerateOTPTokenResponse, error) {
	if err := authorize(ctx, client, constants.MFAGenerate); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	otpRequest := &mfaClient.GenerateOTPTokenRequest{
		Principal:    iam.Credentials.Principal,
		ClientSecret: clientSecret,
	}

	return client.mfaService.GenerateOTPToken(ctx, otpRequest)
}

func (client *Client) ValidateOTP(ctx context.Context, otp string, clientSecret string) (*mfaClient.ValidateOTPTokenResponse, error) {
	if err := authorize(ctx, client, constants.MFAValidate); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	validateRequest := &mfaClient.ValidateOTPTokenRequest{
		Token:        otp,
		Principal:    iam.Credentials.Principal,
		ClientSecret: clientSecret,
	}

	return client.mfaService.ValidateOTPToken(ctx, validateRequest)
}
