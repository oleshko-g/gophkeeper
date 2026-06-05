package grpc

import (
	"crypto/tls"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	DepositorClient = pb.DepositorServiceClient
	KeeperClient    = pb.KeeperServiceClient
	CallOption      = grpc.CallOption
)

func New(gophKeeperURL string) (*Client, error) {
	conn, err := grpc.NewClient(
		gophKeeperURL,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		DepositorClient: pb.NewDepositorServiceClient(conn),
		KeeperClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

type Client struct {
	DepositorClient
	KeeperClient
}
