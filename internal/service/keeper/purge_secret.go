package keeper

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) PurgeSecret(ctx context.Context, in *pb.PurgeSecretRequest) (*pb.PurgeSecretResponse, error) {
	method := "PurgeSecret"
	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}
	err := in.Validate()
	if err != nil {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeRequestValidation, SvcName: s.Name()}, err)
	}

	id, err := uuid.FromBytes(in.SecretId.Bytes)
	if err != nil {
		return nil, service.WrapError(&service.Err{SvcName: s.Name(), Method: method,
			Type: service.ErrTypeRequestValidation},
			err)
	}

	err = s.RemoveSecret(ctx,
		uuidv7.FromString(publicKeyID),
		uuidv7.FromString(id.String()),
	)
	if err != nil {
		return nil, service.WrapError(&service.Err{}, err)
	}

	return &pb.PurgeSecretResponse{}, nil
}
