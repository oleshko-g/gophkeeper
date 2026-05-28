package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oleshko-g/gophkeeper/client/internal/model"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// downloadRunE reads a deposited secret ID from the standard input, downloads the deposited secret data, decrypts it and prints to the standard output
func (app *a) downloadRunE(cmd *cobra.Command, _ []string) error {
	// read the input
	secretIDByteString, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return nil
	}

	// validate the input
	secretID, err := depositedSecretID(secretIDByteString)
	if err != nil {
		return err
	}

	// prepare the ctx
	ctx, cancel := context.WithTimeout(cmd.Context(), 100*time.Millisecond)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", app.config.AuthToken)

	// make the request
	res, err := app.RetrieveSecret(ctx, &pb.RetrieveSecretRequest{
		SecretId: secretID,
	})
	if err != nil {
		return err
	}

	// unmarshal the data into the [model.DepositedSecret]
	secret := &model.Secret{}
	payload := res.GetPayload()
	err = proto.Unmarshal(payload, secret)
	if err != nil {
		return err
	}

	// decrypt the DEK
	decryptedDEK, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, app.privateKey, secret.EncryptedDek, nil)
	if err != nil {
		return err
	}

	// decrypt tht Data
	protoOpenSecretData, err := decryptSecret(secret.EncryptedData, decryptedDEK, secret.Nonce)
	if err != nil {
		return err
	}

	// unmarshal into [model.OpenSecret]
	openSecret := &model.OpenSecret{}
	err = proto.Unmarshal(protoOpenSecretData, openSecret)
	if err != nil {
		return err
	}

	// print the opened secret
	fmt.Fprintf(os.Stdout, "%+v", openSecret)

	return nil
}
