package depositor

import (
	"context"
	"encoding/pem"
	"errors"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

// Register registers an anonymous public key and returns a refresh token.
// The owner of the refresh token can then [Authorize] apps to [Connect] to [KeeperService]
func (s *Service) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if !in.ProtoReflect().IsValid() {
		return nil, errRegisterRequestIsEmpty
	}

	if err := in.Validate(); err != nil {
		return nil, &service.Err{
			SvcName: "Depositor",
			Method:  "Register",
			Err:     err,
		}
	}

	pubKeyID, err := s.Depositor.StorePubKey(ctx, in.GetPubKey())
	if err != nil {
		return nil, err
	}

	rt := depositor.RefreshToken{
		PubKeyID: uuidv7.FromString(pubKeyID),
		ID:       uuidv7.New(),
		TTL:      s.refreshTokenTTL,
	}
	err = s.Depositor.StoreRefreshToken(ctx, rt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func validatePubKey(s string) error {
	pem, _ := pem.Decode([]byte(s))
	if pem == nil {
		return errDecodingPEM
	}
	return nil
}

var (
	errDecodingPEM = errors.New("decoding PEM")
)
