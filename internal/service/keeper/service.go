package keeper

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
)

var _ service.Keeper = (*Service)(nil)

func New(s storage.Keeper) *Service {
	return &Service{name: "Keeper", Keeper: s}
}

type Service struct {
	name string
	storage.Keeper
	pb.UnimplementedKeeperServiceServer
}

func (s *Service) DepositSecret(ctx context.Context, req *pb.DepositSecretRequest) (*pb.DepositSecretResponse, error) {

	publicKeyID, ok := service.PubKeyIDFromCtx(ctx)
	if !ok {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeUnauthenticated, SvcName: s.Name()}, errEmptyPubKeyID)
	}

	sec := secret.FromProto(req)

	depositedSecret, err := s.Keeper.StoreSecret(ctx, uuidv7.FromString(publicKeyID), sec)
	if err != nil {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeStorage, SvcName: s.Name()}, err)
	}

	binDepositedSecretID, err := depositedSecret.ID.Value.MarshalBinary()
	if err != nil {
		return nil, service.WrapError(&service.Err{Type: service.ErrTypeInternal, SvcName: s.Name()}, err)
	}

	return &pb.DepositSecretResponse{DepositedSecretId: &pb.UUID{Bytes: binDepositedSecretID}},
		nil
}

func (s *Service) Name() string {
	return s.name
}
