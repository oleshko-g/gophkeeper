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
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return err
	}

	encodedPubKey := &bytes.Buffer{}
	err = pem.Encode(encodedPubKey, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	if err != nil {
		return err
	}

	res, err := app.client.DepositorServiceClient.Register(
		cmd.Context(),
		&pb.RegisterRequest{
			PubKey: new(encodedPubKey.String()),
		},
	)
	if err != nil {
		return err
	}

	refreshToken := res.GetRefreshToken()
	if refreshToken == "" {
		return fmt.Errorf("refresh token is empty")
	}

	err = os.WriteFile(path.Join(app.cfgDir, "refresh_token.txt"), []byte(refreshToken), 0600)
	if err != nil {
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	encodedPrivKey := &bytes.Buffer{}
	err = pem.Encode(encodedPrivKey, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})
	if err != nil {
		return err
	}

	// write RSA key pair
	err = os.WriteFile(path.Join(app.cfgDir, "id_rsa"), encodedPrivKey.Bytes(), 0600)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(app.cfgDir, "id_rsa.pub"), encodedPubKey.Bytes(), 0644)
	if err != nil {
		return err
	}

	app.config.RegisteredPubKey = map[string]string{
		encodedPubKey.String(): refreshToken,
	}

	app.state = pubKeyRegistered

	return nil
}
