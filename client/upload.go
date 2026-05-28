package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/client/internal/model"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (app *a) uploadRunE(cmd *cobra.Command, args []string) error {
	// TODO: add readWithTimeout
	in, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	o := model.OpenSecret{
		Type: new(model.SecretType_SECRET_TYPE_STRING),
		Name: new("secret"),
		Data: []byte(in),
	}
	err = o.Validate()
	if err != nil {
		return err
	}

	protoOpenSecretData, err := proto.Marshal(&o)
	if err != nil {
		return err
	}

	encryptedOpenSecret, DEK, nonce, err := encryptOpenSecret(protoOpenSecretData)
	if err != nil {
		return err
	}

	encryptedDEK, err := rsa.EncryptOAEP(sha256.New(), rand.Reader,
		app.publicKey, DEK, nil,
	)
	if err != nil {
		return err
	}

	s := model.Secret{
		EncryptedDek:  encryptedDEK,
		EncryptedData: encryptedOpenSecret,
		Nonce:         nonce,
	}

	protoSecretData, err := proto.Marshal(&s)
	if err != nil {
		return err
	}

	ctx := metadata.AppendToOutgoingContext(
		cmd.Context(),
		"authorization", app.config.AuthToken,
	)
	res, err := app.client.DepositSecret(ctx, &pb.DepositSecretRequest{
		Payload: protoSecretData,
	})
	if err != nil {
		return err
	}

	depositedSecretId := res.GetDepositedSecretId()
	// TODO: marshal into byte model.DepositedSecret, WriteFile
	err = app.saveDepositedSecretID(depositedSecretId)
	if err != nil {
		return err
	}
	return nil
}
