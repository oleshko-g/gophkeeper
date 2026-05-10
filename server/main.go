package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	pgsql "github.com/oleshko-g/gophkeeper/internal/db/pgx"
	queries "github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/security"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service/keeper"
	storageDepositor "github.com/oleshko-g/gophkeeper/internal/storage/depositor"
	storageKeeper "github.com/oleshko-g/gophkeeper/internal/storage/keeper"
	"google.golang.org/grpc/credentials"
)

func main() {
	app := server.App{}

	err := app.Configure()
	if err != nil {
		panic(err)
	}

	pgConn, err := pgsql.New(app.DB.Config)
	if err != nil {
		panic(err)
	}
	app.DB.Conn = pgConn

	q := queries.New(app.DB.Conn)
	app.SetStorage(
		storageDepositor.New(q),
		storageKeeper.New(q),
	)

	cert, err := tls.LoadX509KeyPair(security.CertFile, security.KeyFile)
	if err != nil {
		err := security.WriteX509KeyPair()
		if err != nil {
			panic(err)
		}
		cert, err = tls.LoadX509KeyPair(security.CertFile, security.KeyFile)
		if err != nil {
			panic(err)
		}
	}

	priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		panic("cert is not an RSA private key")
	}
	app.SetService(
		depositor.New(app.Storage.Depositor, priv),
		keeper.New(app.Storage.Keeper),
	)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	err = app.SetServer(creds)
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
