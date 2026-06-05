package keeper

import (
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
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

func (s *Service) Name() string {
	return s.name
}
