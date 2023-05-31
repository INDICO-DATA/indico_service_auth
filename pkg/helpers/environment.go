package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/INDICO-INNOVATION/indico_service_auth/pkg/constants"
)

type InnovationCredentials struct {
	Type         string `json:"type"`
	PrivateKeyID string `json:"private_key_id"`
	Principal    string `json:"principal"`
	PrivateKey   string `json:"private_key"`
	AuthUri      string `json:"auth_uri"`
	TokenUri     string `json:"token_uri"`
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

	authServer := constants.AuthServer
	if os.Getenv("AUTH_SERVER") != "" {
		authServer = os.Getenv("AUTH_SERVER")
	}

	if os.Getenv("INSECURE") != "true" {
		authServer = fmt.Sprintf("%s:443", authServer)
	}

	return innovationCredentials, authServer
}
