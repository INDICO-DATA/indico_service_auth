SHELL=/bin/bash

run:
	INSECURE=true ENVIRONMENT=local INNOVATION_CREDENTIALS=innovation.json AUTH_SERVER=localhost:7001 go run ./example/main.go

setup:
	go mod init github.com/INDICO-INNOVATION/indico_service_auth
	go mod tidy