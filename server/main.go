package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	pgsql "github.com/oleshko-g/gophkeeper/internal/db/pgx"
	queries "github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service/keeper"
	storageDepositor "github.com/oleshko-g/gophkeeper/internal/storage/depositor"
	storageKeeper "github.com/oleshko-g/gophkeeper/internal/storage/keeper"
)

func main() {
	app := server.App{}

	err := app.I_Configure()
	if err != nil {
		panic(err)
	}

	pgConn, err := pgsql.New(app.DB.Config)
	if err != nil {
		panic(err)
	}
	app.DB.Conn = pgConn

	q := queries.New(app.DB.Conn)
	app.II_SetStorage(
		storageDepositor.New(q),
		storageKeeper.New(q),
	)

	err = app.III_SetSecurity()
	if err != nil {
		panic(err)
	}

	app.IV_SetService(
		depositor.New(app.Storage.Depositor, app.Security.PrivateKey()),
		keeper.New(app.Storage.Keeper),
	)

	err = app.V_SetServer(*app.Security.Credentials())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = app.Run(ctx)
	fmt.Println(err)

	err = app.Stop(err)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
