package grpc

import (
	"github.com/oleshko-g/gophkeeper/internal/service"
	keeper_v1 "github.com/oleshko-g/gophkeeper/proto/v1"
	"google.golang.org/grpc"
)

func New(keeper service.Keeper, depositor service.Depositor) *Server {
	s := &Server{}

	s.Server = grpc.NewServer()

	s.implemented.keeper.service = keeper
	keeper_v1.RegisterKeeperServiceServer(s.Server, s.implemented.keeper)

	s.implemented.depositor.service = depositor
	keeper_v1.RegisterDepositorServiceServer(s.Server, s.implemented.depositor)

	return s
}

type Server struct {
	*grpc.Server
	implemented struct {
		keeper
		depositor
	}
}

type keeper struct {
	keeper_v1.UnimplementedKeeperServiceServer
	service service.Keeper
}

type depositor struct {
	keeper_v1.UnimplementedDepositorServiceServer
	service service.Depositor
}
