package keeper

import (
	"context"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
	_ "github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/transform"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) ListSecrets(ctx context.Context, r *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}

	secrets, err := s.RetrieveSecrets(ctx, uuidv7.FromString(publicKeyID))
	if err != nil {
		return nil, service.WrapError(
			&service.Err{Type: service.ErrTypeStorage, SvcName: s.Name()},
			err)
	}

	secretIDs := transform.SliceToSlice(secrets, func(s secret.DepositedSecret) *pb.UUID {
		bytes, _ := s.ID.Value.MarshalBinary()
		return &pb.UUID{
			Bytes: bytes,
		}
	})

	return &pb.ListSecretsResponse{
		SecretIds: secretIDs,
	}, nil
}
