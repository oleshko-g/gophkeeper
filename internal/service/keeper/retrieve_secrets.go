package keeper

import (
	"context"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
)

func (s *Service) RetrieveSecret(ctx context.Context, in *pb.RetrieveSecretRequest) (*pb.RetrieveSecretResponse, error) {
	// Authenticate the user
	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}

	// RetrieveSecret from storage
	// s.Keeper.

	// return response
	_ = publicKeyID
	return &pb.RetrieveSecretResponse{}, nil
}
