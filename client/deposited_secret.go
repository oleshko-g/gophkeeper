package main

import (
	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
)

// saveDepositedSecretID saves the ID of a deposited secret to the app's depositedSecrets map.
func (app *a) saveDepositedSecretID(depositedSecretId *pb.UUID) error {
	id, err := uuid.FromBytes(depositedSecretId.Bytes)
	if err != nil {
		return err
	}
	app.depositedSecrets[id.String()] = struct{}{}

	return nil
}
