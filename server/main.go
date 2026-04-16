package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	_ "github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service/keeper"
	storageDepositor "github.com/oleshko-g/gophkeeper/internal/storage/depositor"
	storageKeeper "github.com/oleshko-g/gophkeeper/internal/storage/keeper"
)

func main() {
	app := server.App{}

	err := app.Configure()
	if err != nil {
		panic(err)
	}

	app.SetStorage(
		storageKeeper.New(),
		storageDepositor.New(),
	)

	app.SetService(
		keeper.New(app.Storage.Keeper),
		depositor.New(app.Storage.Depositor),
	)

	app.SetServer()

	ctx := context.Background()
	err = app.Run(ctx)
	fmt.Println(err)

	err = app.Stop(err)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
