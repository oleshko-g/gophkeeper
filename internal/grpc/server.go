package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// New creates a new gRPC server with the gRPC server
func New(keeper service.Keeper, depositor service.Depositor) *Server {
	s := &Server{}

	s.Server = grpc.NewServer()
	reflection.Register(s.Server)

	s.implemented.keeper.service = keeper
	pb.RegisterKeeperServiceServer(s.Server, s.implemented.keeper)

	s.implemented.depositor.service = depositor
	pb.RegisterDepositorServiceServer(s.Server, s.implemented.depositor)

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
func (s *Server) Serve(ctx context.Context, cfg *Config) error {
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		errCh <- s.Server.Serve(lis)
	}()

	fmt.Printf("gRPC server is listening on %s\n", cfg.GRPCAddr)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

type keeper struct {
	pb.UnimplementedKeeperServiceServer
	service service.Keeper
}

type depositor struct {
	pb.UnimplementedDepositorServiceServer
	service service.Depositor
}
