package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

func (app *a) deleteRunE(cmd *cobra.Command, _ []string) error {
	l := app.logger.WithGroup("deleteRunE")
	slog.SetLogLoggerLevel(slog.LevelDebug)
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
		l.Debug(fmt.Sprintf("context deadline %v", deadline))
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
	_, err = app.Client.PurgeSecret(ctx, &pb.PurgeSecretRequest{SecretId: secretID})
	if err != nil {
		return err
	}

	delete(app.depositedSecrets, string(data))

	return nil
}
