package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/db/pgx"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

// App is the struct to hold a gophkeeper server implementation.
//
// In order to set up the app call:
//  1. [App.Configure]
//  2. [App.SetStorage]
//  3. [App.SetService]
//  4. [App.SetServer]
type App struct {
	grpc struct {
		*grpc.Config
		*grpc.Server
	}
	DB struct {
		*pgsql.Config
		*sql.DB
	}
	*storage.Storage
	*service.Service
}

// Configure loads environment variables and processes them into the [App.config] field.
func (a *App) Configure(envFiles ...string) error {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return err
	}

	a.grpc.Config = &grpc.Config{}
	err = envconfig.Process("", a.grpc.Config)
	if err != nil {
		return err
	}

	a.DB.Config = &pgsql.Config{}
	err = envconfig.Process("", a.DB.Config)
	if err != nil {
		return err
	}

	return nil
}

// SetStorage populates the [App.Storage] field with the given [storage.Keeper] and [storage.Depositor] implementations.
func (a *App) SetStorage(keeper storage.Keeper, depositor storage.Depositor) {
	a.Storage = &storage.Storage{
		Keeper:    keeper,
		Depositor: depositor,
	}
}

// SetService populates the [App.Service] field with the given [service.Keeper] and [service.Depositor] implementations.
func (a *App) SetService(keeper service.Keeper, depositor service.Depositor) {
	a.Service = &service.Service{
		Keeper:    keeper,
		Depositor: depositor,
	}
}

// Setup initializes [App.Server] with the [App.Service] implementations. If [App.Service] is nil it panics
func (a *App) SetServer() error {
	if a.Service == nil {
		return errors.New("field Service is nil. SetService() must be called before SetServer")
	}

	a.grpc.Server = grpc.New(
		a.Service.Keeper,
		a.Service.Depositor,
	)

	return nil
}

// Run starts the [App.Server] and blocks until it is stopped or an error occurs.
func (a *App) Run(ctx context.Context) error {
	if a.grpc.Config == nil {
		return errors.New("config is nil. App.Configue() must be called before App.Run()")
	}

	if err := a.grpc.Server.Serve(ctx, a.grpc.Config); err != nil {
		return err
	}

	<-ctx.Done()

	return ctx.Err()
}

func (a *App) Stop(err error) error {
	fmt.Printf("Stopping the app because of %q\n", err)
	return a.grpc.Server.Stop()
}
