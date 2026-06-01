package main

import (
	"context"
	"io"
	_ "os"
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	_ "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

func (app *a) deleteRunE(cmd *cobra.Command, _ []string) error {
	var (
		in  io.Reader = cmd.InOrStdin()
		ctx context.Context
	)

	if ctx = cmd.Context(); ctx != nil {
		ctx = context.Background()
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		timeout := 1 * time.Second
		deadline = time.Now().Add(timeout)
	}

	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	data, err := ReadAllCtx(ctx, in)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errEmptyInput
	}

	secretID, err := depositedSecretID(data)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", app.Authorization.Token)
	_, err = app.client.PurgeSecret(ctx, &pb.PurgeSecretRequest{SecretId: secretID})
	if err != nil {
		return err
	}

	return nil
}
