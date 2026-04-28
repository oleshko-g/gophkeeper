package keeper

import (
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
)

var _ service.Keeper = (*Service)(nil)

func New(s storage.Keeper) *Service {
	return &Service{Keeper: s}
}

type Service struct {
	storage.Keeper
	pb.UnimplementedKeeperServiceServer
}
