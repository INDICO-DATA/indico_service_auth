package main

import (
	"context"
	"fmt"
	"log"

	indicoserviceauth "github.com/INDICO-INNOVATION/indico_service_auth"
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
	response, err := client.ListResources(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)

	fmt.Println("")

	response2, err := client.CreateClient(ctx, "Anthony", "user")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response2)

	fmt.Println("")

	response3, err := client.CreateResource(ctx, "use", "Use stuff")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response3)

	fmt.Println("")

	response4, err := client.CreateResourceScope(ctx, "Developer", "dev", "He program", 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response4)

	fmt.Println("")

	response5, err := client.CreateServiceAccount(ctx, "Tony", "tony", "he likes food")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response5)

	fmt.Println("")

	response6, err := client.GetResourceScope(ctx, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response6)
}
