package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
)

func (app *a) authorizeRunE(cmd *cobra.Command, _ []string) error {

	ciphertext, err := base64.RawStdEncoding.DecodeString(app.config.RegisteredPubKeyID)
	if err != nil {
		return err
	}

	data, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		app.PrivateKey,
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
