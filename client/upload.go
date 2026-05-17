package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io"
	"os"
	"path"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/client/internal/model"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (app *a) uploadRunE(cmd *cobra.Command, args []string) error {
	in, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	o := model.OpenSecret{
		Type: new(model.SecretType_SECRET_TYPE_STRING),
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

	encryptedOpenSecret, DEK, nonce, err := encryptOpenSecret(protoOpenSecretData)
	if err != nil {
		return err
	}

	encryptedDEK, err := rsa.EncryptOAEP(sha256.New(), rand.Reader,
		app.config.PublicKey, DEK, nil,
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
	err = app.storeDepositedSecretID(depositedSecretId)
	if err != nil {
		return err
	}
	return nil
}

func encryptOpenSecret(protoOpenSecretData []byte) (encryptedData, DEK, nonce []byte, err error) {
	const (
		keyLen   int = 32
		nonceLen int = 12
	)

	buf := make([]byte, keyLen+nonceLen) // 32 bytes to cipher with AES 256 bit + 12c

	rand.Read(buf) // always reads up the len(b) and never errors

	DEK = buf[:keyLen]
	nonce = buf[keyLen:]

	cipherBlock, err := aes.NewCipher(DEK)
	if err != nil {
		return nil, nil, nil, err
	}

	awed, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedData = awed.Seal(nil, nonce, protoOpenSecretData, nil)

	return encryptedData, DEK, nonce, nil
}

func (app *a) storeDepositedSecretID(depositedSecretId *pb.UUID) error {
	const storageFileName string = "depositedSecrets.json"
	storageFilePath := path.Join(app.cfgDir, storageFileName)
	storageFile, err := os.OpenFile(storageFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0x600)
	if err != nil {
		return err
	}
	defer storageFile.Close()

	id := string(depositedSecretId.Bytes)

	err = json.NewEncoder(storageFile).Encode(model.DepositedSecret{ID: id})
	if err != nil {
		return err
	}

	return nil
}
