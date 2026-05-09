package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/spf13/cobra"
)

// registerRunE registers the user and gets the refresh token
// It generates an RSA key, calls the gophkeeper server to register the user with the public RSA key.
// If successful, it stores the refresh token on disk and sets [app.state] to [pubKeyRegistered]
func (app *a) registerRunE() func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {

		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}

		pubKey, err := encodeRSAPubKey(priv)
		if err != nil {
			return err
		}
		_ = pubKey

		return nil
	}
}

func encodeRSAPubKey(priv *rsa.PrivateKey) ([]byte, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	pubKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return pubKey, nil
}
