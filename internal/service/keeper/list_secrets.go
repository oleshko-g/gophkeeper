package keeper

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
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

	secretIDs, err := s.Keeper.RetrieveSecretIDs(ctx, uuidv7.FromString(publicKeyID))
	if err != nil {
		return nil, service.WrapError(
			&service.Err{Type: service.ErrTypeStorage, SvcName: s.Name()},
			err)
	}

	secretIDsBytes := transform.SliceToSlice(secretIDs, func(s uuid.UUID) *pb.UUID {
		bytes, _ := s.MarshalBinary()
		return &pb.UUID{
			Bytes: bytes,
		}
	})

	return &pb.ListSecretsResponse{
		SecretIds: secretIDsBytes,
	}, nil
}
