package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
)

func (app *a) authorizeRunE(cmd *cobra.Command, _ []string) error {
	pemBytes, err := os.ReadFile(app.config.PrivateKeyFilePath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return errors.New("failed to decode PEM")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	RSAprivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("not RSA private key")
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(app.config.RegisteredPubKeyID)
	if err != nil {
		return err
	}

	data, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		RSAprivKey,
		ciphertext,
		nil)
	if err != nil {
		return err
	}

	res, err := app.client.Authorize(
		cmd.Context(),
		&pb.AuthorizeRequest{
			DecryptedId: new(string(data)),
		})
	if err != nil {
		return err
	}

	app.config.AuthToken = res.GetAuthToken()

	return nil
}
