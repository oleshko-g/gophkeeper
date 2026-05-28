package main

import (
	"strings"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
)

// saveDepositedSecretID saves the ID of a deposited secret to the app's depositedSecrets map.
func (app *a) saveDepositedSecretID(depositedSecretId *pb.UUID) error {
	l := app.logger.WithGroup("saveDepositedSecretID")
	id, err := uuid.FromBytes(depositedSecretId.Bytes)
	if err != nil {
		l.Error(err.Error())
		return err
	}
	app.depositedSecrets[id.String()] = struct{}{}

	return nil
}

func depositedSecretID(uuidByteString []byte) (*pb.UUID, error) {
	l := app.logger.WithGroup("depositedSecretID")

	s, _ := strings.CutSuffix(string(uuidByteString), "\n")

	id, err := uuid.Parse(s)
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	idBytes, err := id.MarshalBinary()
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	return &pb.UUID{
		Bytes: idBytes,
	}, nil
}
