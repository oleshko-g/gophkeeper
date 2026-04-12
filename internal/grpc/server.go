package grpc

import (
	keeper_v1 "github.com/oleshko-g/gophkeeper/proto/v1"
)

func New() *Server {

	return &Server{}
}

type Server struct {
	keeper_v1.UnimplementedClientAppServiceServer
	keeper_v1.UnimplementedKeeperServiceServer
	*implemented
}

type implemented struct {}
