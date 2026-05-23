package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
)

// listRunE requests the deposited secrets and outputs their IDs to an [io.Writer]
func (app *a) listRunE(cmd cobra.Command, _ []string) error {
	l := app.logger.WithGroup("listRunE")

	ctx, cancel := context.WithTimeout(cmd.Context(), 100*time.Millisecond)
	defer cancel()

	depositedSecretIDs, err := app.client.ListSecrets(ctx, &pb.ListSecretsRequest{})
	if err != nil {
		l.Error("list secrets", "error", err)
		return err
	}

	for i, secretID := range depositedSecretIDs.SecretIds {
		if secretID == nil {
			l.Warn("nil secret ID", "id number", i+1)
			continue
		}

		id, err := uuid.FromBytes(secretID.Bytes)
		if err != nil {
			l.Error("unmarshalling secret ID bytes", "id number", i+1)
			return err
		}
		app.depositedSecrets[id.String()] = struct{}{}
	}

	return nil
}
