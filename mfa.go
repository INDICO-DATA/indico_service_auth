package indicoserviceauth

// func (client *Client) GenerateOTP(ctx context.Context) (*mfaClient.GenerateOTPTokenResponse, error) {
// 	if err := authorize(ctx, client, "mfa.validate"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	otpRequest := &mfaClient.GenerateOTPTokenRequest{
// 		ClientId:     iam.Credentials.ClientID,
// 		ClientSecret: iam.Credentials.ClientSecret,
// 	}

// 	return client.mfaService.GenerateOTPToken(ctx, otpRequest)
// }

// func (client *Client) ValidateOTP(ctx context.Context, otp string, useSecret bool) (*mfaClient.ValidateOTPTokenResponse, error) {
// 	if err := authorize(ctx, client, "mfa.validate"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	validateRequest := &mfaClient.ValidateOTPTokenRequest{
// 		Token:    otp,
// 		ClientId: iam.Credentials.ClientID,
// 	}

// 	if useSecret {
// 		validateRequest.ClientSecret = iam.Credentials.ClientSecret
// 	}

// 	return client.mfaService.ValidateOTPToken(ctx, validateRequest)
// }

// func (client *Client) GenerateSecretKey(ctx context.Context) (*mfaClient.OTPSecretResponse, error) {
// 	if err := authorize(ctx, client, "mfa.validate"); err != nil {
// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	secretRequest := &mfaClient.GenerateOTPTokenRequest{
// 		ClientId: iam.Credentials.ClientSecret,
// 	}

// 	return client.mfaService.GenerateSecretKey(ctx, secretRequest)
// }
