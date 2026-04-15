package server

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

func Build(
	grpcConfig *grpc.Config,

	keeperStorage storage.Keeper,
	keeperService service.Keeper,

	depositorStorage storage.Depositor,
	depositorService service.Depositor,
) (*app, error) {
	a := &app{}

	a.grpc.Config = *grpcConfig

	a.keeper.storage = keeperStorage
	a.keeper.service = keeperService

	a.depositor.storage = depositorStorage
	a.depositor.service = depositorService

	a.grpc.Server = grpc.New(a.keeper.service, a.depositor.service)

	return a, nil
}

type app struct {
	grpc struct {
		grpc.Config
		Server *grpc.Server
	}
	depositor struct {
		service service.Depositor
		storage storage.Depositor
	}
	keeper struct {
		service service.Keeper
		storage storage.Keeper
	}
}

// configue initializes the gRPC config from environment variables.
func Configue(grpcConfig *grpc.Config) error {
	return envconfig.Process("", grpcConfig)
}

func (a *app) Run() error {
	if a.grpc.Server == nil {
		return ErrServerIsNil
	}

	return a.grpc.Server.Serve(&a.grpc.Config)
}

func (a *app) Stop() error {
	return a.grpc.Server.Stop()
}
