package depositor

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
)

// Register registers an anonymous RSA public key and returns the ID encrypted with the registered key.
func (s *Service) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	methodName := "Register"

	if !in.ProtoReflect().IsValid() {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, errRequestEmpty)
	}

	if err := in.Validate(); err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, err)
	}

	depositorPubKey, err := parseRSAPubKey(in.GetPubKey())
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, err)
	}

	pubKeyID, err := s.Depositor.StorePubKey(ctx, in.GetPubKey())
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeStorage, err)
	}

	rsaEncryptedID, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		depositorPubKey,
		[]byte(pubKeyID),
		nil,
	)
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeInternal, err)
	}

	rsaEncryptedIDBase64Str := base64.RawStdEncoding.EncodeToString(rsaEncryptedID)

	return &pb.RegisterResponse{
			EncryptedId: &rsaEncryptedIDBase64Str,
		},
		nil
}

// parseRSAPubKey parses an RSA public key from a PEM-encoded string.
func parseRSAPubKey(s string) (key *rsa.PublicKey, err error) {
	pem, _ := pem.Decode([]byte(s))
	if pem == nil {
		return nil, errDecodingPEM
	}

	key, err = x509.ParsePKCS1PublicKey(pem.Bytes)
	if err == nil {
		return key, nil
	}

	pub, err := x509.ParsePKIXPublicKey(pem.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errNotRSA
	}

	return key, nil
}

var (
	errDecodingPEM = errors.New("decoding PEM")
	errNotRSA      = errors.New("not an RSA pub key")
)
