package main

import (
	"context"
	"fmt"
	"log"

	indicoserviceauth "github.com/indicoinnovation/indico_service_auth"
)

func main() {
	testMfa()
}

func testMfa() {
	client, ctx, err := indicoserviceauth.NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("MFA Client successfully connected to IAM through Indico Service Auth")

	testGenerateAndValidate(client, ctx)
}

func testGenerateAndValidate(client *indicoserviceauth.Client, ctx context.Context) {
	clientSecret := "INNOVATION"
	response, err := client.GenerateOTP(ctx, clientSecret)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Generate OTP Response:")
	fmt.Printf("%+v\n\n", response)

	responsev, err := client.ValidateOTP(ctx, response.Token, clientSecret)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Validate OTP Response:")
	fmt.Printf("%+v\n\n", responsev)
}
