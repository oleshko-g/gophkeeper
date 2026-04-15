package server

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service"
)

func Build(
	keeperService service.Keeper,
	depositorService service.Depositor,
) (*app, error) {
	a := &app{}

	envconfig.Process("", a.GRPC.Config)

	a.keeper.service = keeperService
	a.depositor.service = depositorService

	a.GRPC.Server = grpc.New(a.keeper.service, a.depositor.service)

	return a, nil
}

type app struct {
	GRPC struct {
		grpc.Config
		Server *grpc.Server
	}
	depositor struct {
		service service.Depositor
	}
	keeper struct {
		service service.Keeper
	}
}

func (a *app) Run() error {
	if a.GRPC.Server == nil {
		return ErrServerIsNil
	}

	return a.GRPC.Server.Serve(&a.GRPC.Config)
}

func (a *app) Stop() error {
	return a.GRPC.Server.Stop()
}
