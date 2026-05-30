package main

import (
	"context"
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

// listRunE requests the deposited secrets and outputs their IDs to an [io.Writer]
func (app *a) listRunE(cmd *cobra.Command, _ []string) error {
	l := app.logger.WithGroup("listRunE")

	ctx, cancel := context.WithTimeout(cmd.Context(), 100*time.Millisecond)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", app.Authorization.Token)
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

		err := app.saveDepositedSecretID(secretID)
		if err != nil {
			l.Error("unmarshalling secret ID bytes", "id number", i+1)
			return err
		}
	}

	return nil
}
