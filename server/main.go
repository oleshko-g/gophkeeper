package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service/keeper"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

type app struct {
	*grpc.Config
	*grpc.Server
	depositor struct {
		*depositor.Service
		depositor.Storage
	}
	keeper struct {
		*keeper.Service
		keeper.Storage
	}
}

func (a *app) configure() error {
	return envconfig.Process("", a.Config)
}

func (a *app) setup() error {
	err := a.configure()
	if err != nil {
		return err
	}

	a.keeper.Storage = storage.NewKeeper()
	a.depositor.Storage = storage.NewDepositor()

	a.keeper.Service = keeper.New(a.keeper.Storage)
	a.depositor.Service = depositor.New(a.depositor.Storage)

	a.Server = grpc.New(a.keeper.Service, a.depositor.Service)

	return nil
}

func (a *app) run() error {
	return a.Server.Serve(a.Config)
}

func (a *app) stop() error {
	return nil
}

func main() {
	a := &app{
		Config: &grpc.Config{},
	}

	godotenv.Load(".env")
	if err := a.configure(); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", a.Config)

	if err := a.setup(); err != nil {
		panic(err)
	}
	if err := a.run(); err != nil {
		panic(err)
	}
}
