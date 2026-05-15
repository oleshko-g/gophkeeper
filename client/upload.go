package main

import (
	"io"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

func (app *a) uploadRunE(cmd *cobra.Command, args []string) error {
	in, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	o := OpenSecret{
		Type: new(SecretType_SECRET_TYPE_STRING),
		Name: new("secret"),
		Data: []byte(in),
	}
	err = o.ValidateAll()
	if err != nil {
		return err
	}

	protoOpenSecretData, err := proto.Marshal(&o)
	if err != nil {
		return err
	}

	// TODO:
	encryptedOpenSecret, encryptedDEK, nonce, err := encryptOpenSecret(protoOpenSecretData)
	if err != nil {
		return err
	}

	s := Secret{
		EncryptedDek:  encryptedDEK,
		EncryptedData: encryptedOpenSecret,
		Nonce:         nonce,
	}

	protoSecretData, err := proto.Marshal(&s)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	res, err := app.client.DepositSecret(ctx, &pb.DepositSecretRequest{
		Payload: protoSecretData,
	})
	if err != nil {
		return err
	}

	_ = protoSecretData
	_ = res
	return nil
}

func encryptOpenSecret(protoSecretData []byte) (encryptedData, dek, nonce []byte, err error) {
	// TODO:
	// * encrypt protoSecretData with AES-GCM
	// * encrypt dek
	return nil, nil, nil, nil
}
