package depositor

import (
	"context"
	"errors"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) Authorize(ctx context.Context, in *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	methodName := "Authorize"

	pk, err := s.RetrievePubKeyByID(ctx, uuidv7.FromString(in.GetDecryptedId()))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, s.wrapError(methodName, service.ErrTypeBusinessLogic, err)
		}
		return nil, s.wrapError(methodName, service.ErrTypeStorage, err)
	}

	// issue the JWT with claims: ID, App.UUIDv7
	// encode JWT to Base64
	// return as Token
	_ = pk
	return &pb.AuthorizeResponse{}, nil
}
