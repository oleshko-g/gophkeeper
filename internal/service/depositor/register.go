package depositor

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

// Register registers an anonymous RSA public key and returns a refresh token.
// The owner of the refresh token can then [Authorize] apps to [Connect] to [KeeperService]
func (s *Service) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	methodName := "Register"

	if !in.ProtoReflect().IsValid() {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, errRequestEmpty)
	}

	if err := in.Validate(); err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, err)
	}

	if err := validateRSAPubKey(in.GetPubKey()); err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeRequestValidation, err)
	}

	pubKeyID, err := s.Depositor.StorePubKey(ctx, in.GetPubKey())
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeStorage, err)
	}

	rt := depositor.RefreshToken{
		PubKeyID: uuidv7.FromString(pubKeyID),
		ID:       uuidv7.New(),
		TTL:      s.refreshTokenTTL,
	}
	err = s.Depositor.StoreRefreshToken(ctx, rt)
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeStorage, err)
	}

	return &pb.RegisterResponse{RefreshToken: &rt.ID.String}, nil
}

func validateRSAPubKey(s string) error {
	pem, _ := pem.Decode([]byte(s))
	if pem == nil {
		return errDecodingPEM
	}

	_, err := x509.ParsePKCS1PublicKey(pem.Bytes)
	if err == nil {
		return nil
	}

	pub, err := x509.ParsePKIXPublicKey(pem.Bytes)
	if err != nil {
		return err
	}

	if _, ok := pub.(*rsa.PublicKey); !ok {
		return errNotRSA
	}

	return nil
}

var (
	errDecodingPEM = errors.New("decoding PEM")
	errNotRSA      = errors.New("not an RSA pub key")
)
