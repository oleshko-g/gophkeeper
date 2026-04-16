package grpc

import (
	"net"

	"github.com/oleshko-g/gophkeeper/internal/service"
	keeper_v1 "github.com/oleshko-g/gophkeeper/proto/v1"
	"google.golang.org/grpc"
)

// New creates a new gRPC server with the gRPC server
func New(keeper service.Keeper, depositor service.Depositor) *Server {
	s := &Server{}

	s.Server = grpc.NewServer()

	s.implemented.keeper.service = keeper
	keeper_v1.RegisterKeeperServiceServer(s.Server, s.implemented.keeper)

	s.implemented.depositor.service = depositor
	keeper_v1.RegisterDepositorServiceServer(s.Server, s.implemented.depositor)

	return s
}

// Server is a gRPC server with the implemented [minifier_v1.Keeper] and [minifier_v1.Depositor] services
type Server struct {
	*grpc.Server
	implemented struct {
		keeper
		depositor
	}
}

// Stop stops the gRPC server gracefully.
func (s *Server) Stop() error {
	s.Server.GracefulStop()
	return nil
}

// Serve starts the gRPC server on the configured address and returns any error that occurs.
func (s *Server) Serve(cfg *Config) error {
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return err
	}

	return s.Server.Serve(lis)
}

type keeper struct {
	keeper_v1.UnimplementedKeeperServiceServer
	service service.Keeper
}

type depositor struct {
	keeper_v1.UnimplementedDepositorServiceServer
	service service.Depositor
}
