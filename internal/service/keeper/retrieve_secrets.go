package keeper

import (
	"context"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) RetrieveSecret(ctx context.Context, in *pb.RetrieveSecretRequest) (*pb.RetrieveSecretResponse, error) {
	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}

	secretData, err := s.Keeper.RetrieveSecret(ctx, uuidv7.FromString(publicKeyID))
	if err != nil {
		return nil, service.WrapError(
			&service.Err{Type: service.ErrTypeStorage, SvcName: s.Name()},
			err)
	}

	return &pb.RetrieveSecretResponse{Payload: secretData}, nil
}
