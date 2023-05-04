# Indico Innovation MFA Client

[![pkg-img]][pkg-url]

## Rationale

We were looking for a simple way to use Indico Innovation MFA service in Indico Innovation IAM GRPC server. We will have to develop a new GRPC client every time a client application need to connect to IAM server. 

So we have built this MFA Client to be imported by any Golang application instead to deploy it as a microservice avoiding to make http requests to work.

## Dependencies

* Auth, MFA protobuffers
* Environment variables:
    > INNOVATION_CREDENTIALS=<path_to_innovation.json>

    > AUTH_SERVER=localhost:50051

## Features

* Connect to MFA Service in IAM GRPC server.

## Install
Go version 1.14+
```
go get github.com/INDICO-INNOVATION/indico_service_auth
```


```go
# Example
import indicoserviceauth "github.com/INDICO-INNOVATION/indico_service_auth"

client, ctx, err := indicoserviceauth.NewClient()
if err != nil {
    log.Fatalf(err.Error())
}

response, err := client.GenerateOTP(ctx)
if err != nil {
    log.Fatalf(err.Error())
}

fmt.Println("Generate OTP Response:")
fmt.Printf("%+v\n", response)

responsev, err := client.ValidateOTP(ctx, response.Token, true)
if err != nil {
    log.Fatalf(err.Error())
}

fmt.Println("Validate OTP Response:")
fmt.Printf("%+v\n", responsev)

ctx.Done()
```

## Documentation

See in [docs](https://pkg.go.dev/github.com/INDICO-INNOVATION/indico_auth_service).

## License

[GNU General Public License v3.0](./LICENSE).
