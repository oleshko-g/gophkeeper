package keeper

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) RetrieveSecret(ctx context.Context, in *pb.RetrieveSecretRequest) (*pb.RetrieveSecretResponse, error) {
	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}
	err := in.Validate()
	if err != nil {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeRequestValidation, SvcName: s.Name()}, err)
	}

	id, err := uuid.FromBytes(in.GetSecretId().Bytes)
	if err != nil {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeRequestValidation, SvcName: s.Name()}, err)
	}

	secretData, err := s.Keeper.RetrieveSecret(ctx,
		uuidv7.FromString(publicKeyID),
		uuidv7.FromString(id.String()),
	)
	if err != nil {
		if errors.Is(err, storage.ErrOwnerMismatch) {
			return nil, service.WrapError(
				&service.Err{Type: service.ErrTypeBusinessLogic, SvcName: s.Name()}, err)
		}
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeStorage, SvcName: s.Name()}, err)
	}

	return &pb.RetrieveSecretResponse{Payload: secretData}, nil
}
