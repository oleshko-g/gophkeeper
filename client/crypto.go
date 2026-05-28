package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

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

func decryptSecret(protoSecretData, DEK, nonce []byte) (protoOpenSecretData []byte, err error) {
	cipherBlock, err := aes.NewCipher(DEK)
	if err != nil {
		return nil, err
	}

	AWED, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	protoOpenSecretData, err = AWED.Open(protoOpenSecretData, nonce, protoSecretData, nil)
	if err != nil {
		return nil, err
	}

	return protoOpenSecretData, nil
}
