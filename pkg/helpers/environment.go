package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type InnovationCredentials struct {
	Type         string `json:"type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	AuthURI      string `json:"auth_uri"`
}

func ParseEnvironment() (*InnovationCredentials, string) {
	var innovationCredentials *InnovationCredentials

	if os.Getenv("INNOVATION_CREDENTIALS") != "" {
		filebyte, err := ioutil.ReadFile(os.Getenv("INNOVATION_CREDENTIALS"))
		if err != nil {
			log.Fatalf("could not find INNOVATION_CREDENTIALS file: %s", err.Error())
		}

		if err = json.Unmarshal([]byte(filebyte), &innovationCredentials); err != nil {
			log.Fatalf("could not parse INNOVATION_CREDENTIALS environment: %s", err.Error())
		}
	} else {
		log.Println("environment variable INNOVATION_CREDENTIALS not set")
	}

	const AUTH_SERVER = "iam.services.indicoinnovation.pt"

	var authServer string = AUTH_SERVER
	if os.Getenv("AUTH_SERVER") != "" {
		authServer = os.Getenv("AUTH_SERVER")
	}

	if os.Getenv("ENVIRONMENT") != "local" {
		authServer = fmt.Sprintf("%s:443", authServer)
	}

	return innovationCredentials, authServer
}
