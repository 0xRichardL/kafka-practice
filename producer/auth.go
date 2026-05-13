package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
)

var (
	priv ed25519.PrivateKey
)

func getSigner() ed25519.PrivateKey {
	if priv == nil {
		priv = ed25519.NewKeyFromSeed([]byte(os.Getenv("PRIVATE_KEY")))
	}
	return priv
}

func SignData(data ...any) string {
	content := ""
	for _, value := range data {
		content += fmt.Sprint(value)
	}
	signature := ed25519.Sign(getSigner(), []byte(content))
	return base64.StdEncoding.EncodeToString(signature)
}
