package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 131072

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyFile, err := os.Create("cert/private.key")
	if err != nil {
		fmt.Println(err)

		return
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Println(err)

		return
	}

	var privateKeyPEM = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		fmt.Println(err)

		return
	}

	publicKey := privateKey.PublicKey

	publicKeyFile, err := os.Create("cert/public.key")
	if err != nil {
		fmt.Println(err)

		return
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println(err)

		return
	}

	var publicKeyPEM = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		fmt.Println(err)

		return
	}
}
