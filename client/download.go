package main

import (
	"context"
	"io"
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
	depositedSecret := &model.DepositedSecret{}
	payload := res.GetPayload()
	err = proto.Unmarshal(payload, depositedSecret)
	if err != nil {
		return err
	}

	// decrypt the DEK

	// decrypt tht Data

	// unmarshal into [model.OpenSecret]

	// print the opened secret

	_ = secretID
	_ = res
	return nil
}
