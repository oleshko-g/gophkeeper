package client

import (
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(cfg *Config) (*Depositor, error) {
	conn, err := grpc.NewClient(
		cfg.GophKeeperURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Depositor{
		DepositorServiceClient: pb.NewDepositorServiceClient(conn),
		KeeperServiceClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

type Depositor struct {
	pb.DepositorServiceClient
	pb.KeeperServiceClient
}
