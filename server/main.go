package main

import (
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service"
)

type app struct {
	*grpc.Config
	*grpc.Server
	service.Depositor
	service.Keeper
}

func (a *app) setup() error {
	return nil
}

func main() {
}
