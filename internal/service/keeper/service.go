//go:build ignore
package keeper

import (
	"context"

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
}

// UploadData uploads data of the owners of the [Register]ed public keys.
func (s *Service) UploadData(ctx context.Context, in *pb.UploadDataRequest) (*pb.UploadDataResponse, error) {
	return nil, nil
}

// DownloadData downloads the [Upload]ed data
func (s *Service) DownloadData(ctx context.Context, in *pb.DownloadDataRequest) (*pb.DownloadDataResponse, error) {
	return nil, nil
}

// DownloadData lists the [Upload]ed data
func (s *Service) ListData(ctx context.Context, in *pb.ListDataRequest) (*pb.ListDataResponse, error) {
	return nil, nil
}

// Delete deletes the [Upload]ed data from the [KeeperService]
func (s *Service) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	return nil, nil
}
