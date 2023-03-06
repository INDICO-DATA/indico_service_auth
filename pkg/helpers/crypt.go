package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParsePublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	public := []byte(publicKeyString)
	publicPem, _ := pem.Decode(public)
	publicKey, err := x509.ParsePKIXPublicKey(publicPem.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func ParsePublicFromPrivateKey(private string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(private))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("NOT PRIVATE KEY: " + block.Type)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	parsedPrivateKey := key.(*rsa.PrivateKey)

	publicKeyDer, err := x509.MarshalPKIXPublicKey(&parsedPrivateKey.PublicKey)
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	pubKeyPem := string(pem.EncodeToMemory(&pubKeyBlock))

	return ParsePublicKey(pubKeyPem)
}
