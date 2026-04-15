package service

import (
	"context"
)

// KeeperService is the interface to:
//   - [UploadData] of the owners of the [Register]ed public keys
//   - [DownloadData] the uploaded data
//   - [ListData] the uploaded data
//   - [DeleteData] the uploaded data
//
//go:generate moq -out keeper_mock.go . Keeper
type Keeper interface {
	// UploadData uploads data of the owners of the [Register]ed public keys.
	UploadData(ctx context.Context, in *UploadDataRequest) (*UploadDataResponse, error)
	// DownloadData downloads the [Upload]ed data
	DownloadData(ctx context.Context, in *DownloadDataRequest) (*DownloadDataResponse, error)
	// DownloadData lists the [Upload]ed data
	ListData(ctx context.Context, in *ListDataRequest) (*ListDataResponse, error)
	// Delete deletes the [Upload]ed data from the [KeeperService]
	DeleteData(ctx context.Context, in *DeleteDataRequest) (*DeleteDataResponse, error)
}

type UploadDataRequest struct {
}

type UploadDataResponse struct {
}

type DownloadDataRequest struct {
}

type DownloadDataResponse struct {
}

type ListDataRequest struct {
}

type ListDataResponse struct {
}

type DeleteDataRequest struct {
}

type DeleteDataResponse struct {
}
