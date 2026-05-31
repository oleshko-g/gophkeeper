package main

import (
	"context"
	"io"
	"os"
	_ "os"
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	_ "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
	"github.com/spf13/cobra"
)

func (app *a) deleteRunE(cmd *cobra.Command, _ []string) error {
	var in io.Reader = cmd.InOrStdin()

	deadline, ok := cmd.Context().Deadline()
	if !ok {
		timeout := 1 * time.Second
		deadline = time.Now().Add(timeout)
	}

	if v, ok := in.(*os.File); ok {
		if err := v.SetReadDeadline(deadline); err != nil {
			return err
		}
	}

	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errEmptyInput
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 150*time.Millisecond)
	defer cancel()

	secretID, err := depositedSecretID(data)
	if err != nil {
		return err
	}

	_, err = app.client.PurgeSecret(ctx, &pb.PurgeSecretRequest{SecretId: secretID})
	if err != nil {
		return err
	}

	return nil
}
