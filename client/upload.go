package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

func encryptOpenSecret(protoSecretData []byte) (encryptedData, DEK, nonce []byte, err error) {
	buf := make([]byte, 32) // 32 bytes to cipher with AES 256 bit
	// make DEK
	b := bytes.NewBuffer(buf)
	rand.Read(b.Bytes()) // always reads up the len(b)
	DEK = b.Bytes()

	c, err := aes.NewCipher(DEK)
	if err != nil {
		return nil, nil, nil, err
	}

	awed, err := cipher.NewGCM(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// make nonce
	b.Reset()
	b.Grow(awed.NonceSize())
	rand.Read(b.Bytes()) // always reads up the len(b)
	nonce = b.Bytes()

	encryptedData = awed.Seal(protoSecretData[:0], nonce, protoSecretData, nil)

	return encryptedData, DEK, nonce, nil
}
