package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	_ "github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service/keeper"
)

type app struct {
	*grpc.Config
	*grpc.Server
	depositor struct {
		depositor.Service
		depositor.Storage
	}
	keeper struct {
		keeper.Service
		keeper.Storage
	}
}

func (a *app) setup() error {
	// 0. configure the gRPC server listener

	// 1. create storage

	// 2. create services

	// 3. create gRPC service server implementations

	// 4. Register gRPC service server implementations

	// 5. Create the gRPC server

	return nil
}

func (a *app) configure() error {
	if err := envconfig.Process("", a.Config); err != nil {
		return err
	}

	return nil
}

func (a *app) run() error {
	return nil
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
