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
	sk := storageKeeper.NewKeeper()
	sd := storageDepositor.NewDepositor()
	a, err := server.Build(
		keeper.New(sk),
		depositor.New(sd),
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", a.GRPC.Config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := a.Run(); err != nil {
			cancel()
		}
	}()
	fmt.Printf("gophkeeper is listening on %s", a.GRPC.Config.GRPCAddr)

	<-ctx.Done()
	err = a.Stop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
