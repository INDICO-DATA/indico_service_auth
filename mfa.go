package go_mfaservice

import (
	mfaClient "github.com/INDICO-INNOVATION/indico_service_auth/client/mfa"
	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/helpers"
	grpcHelper "github.com/INDICO-INNOVATION/indico_service_auth/pkg/helpers/grpc"
)

var mfaservice = mfaClient.NewMFAServiceClient(grpcHelper.Connect("localhost:50051"))

func GenerateOTP(clientID string, clientSecret string) (*mfaClient.GenerateTOTPTokenResponse, error) {
	context, cancel := helpers.InitContext()
	defer cancel()

	otpRequest := &mfaClient.GenerateTOTPTokenRequest{
		ClientId:     clientID,
		ClientSecret: clientSecret,
	}

	return mfaservice.GenerateTOTPToken(context, otpRequest)
}

func ValidateOTP(token int32, clientID string, clientSecret string) (*mfaClient.ValidateTOTPTokenResponse, error) {
	context, cancel := helpers.InitContext()
	defer cancel()

	validateRequest := &mfaClient.ValidateTOTPTokenRequest{
		Token:        token,
		ClientId:     clientID,
		ClientSecret: clientSecret,
	}

	return mfaservice.ValidateTOTPToken(context, validateRequest)
}

func GenerateSecretKey(clientID string) (*mfaClient.TOTPSecretResponse, error) {
	context, cancel := helpers.InitContext()
	defer cancel()

	secretRequest := &mfaClient.GenerateTOTPTokenRequest{
		ClientId: clientID,
	}

	return mfaservice.GenerateSecretKey(context, secretRequest)
}
