package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
)

// registerRunE registers the user and gets the refresh token
// It generates an RSA key, calls the gophkeeper server to register the user with the public RSA key.
// If successful, it stores the refresh token on disk and sets [app.state] to [pubKeyRegistered]
func (app *a) registerRunE(cmd *cobra.Command, _ []string) error {
	l := app.logger.WithGroup("registerRunE")
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		l.Error(err.Error(), "func", "rsa.GenerateKey")
		return err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		l.Error(err.Error(), "func", "x509.MarshalPKIXPublicKey")
		return err
	}

	encodedPubKey := &bytes.Buffer{}
	err = pem.Encode(encodedPubKey, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	if err != nil {
		l.Error(err.Error(), "func", "pem.Encode")
		return err
	}

	res, err := app.client.DepositorServiceClient.Register(
		cmd.Context(),
		&pb.RegisterRequest{
			PubKey: new(encodedPubKey.String()),
		},
	)
	if err != nil {
		l.Error(err.Error(), "func", "app.client.DepositorServiceClient.Register")
		return err
	}

	encryptedID := res.GetEncryptedId()
	if encryptedID == "" {
		err = fmt.Errorf("encryptedID is empty")
		l.Error(err.Error(), "func", "res.GetEncryptedId")
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		l.Error(err.Error(), "func", "x509.MarshalPKCS8PrivateKey")
		return err
	}
	encodedPrivKey := &bytes.Buffer{}
	err = pem.Encode(encodedPrivKey, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})
	if err != nil {
		l.Error(err.Error(), "func", "pem.Encode")
		return err
	}

	// write RSA key pair
	app.config.PrivateKeyFilePath = path.Join(app.cfgDir, "id_rsa")
	err = os.WriteFile(app.config.PrivateKeyFilePath, encodedPrivKey.Bytes(), 0600)
	if err != nil {
		l.Error(err.Error(), "func", "os.WriteFile")
		return err
	}

	app.config.PublicKeyFilePath = path.Join(app.cfgDir, "id_rsa.pub")
	err = os.WriteFile(app.config.PublicKeyFilePath, encodedPubKey.Bytes(), 0644)
	if err != nil {
		l.Error(err.Error(), "func", "os.WriteFile")
		return err
	}

	app.config.RegisteredPubKeyID = res.GetEncryptedId()

	return nil
}
