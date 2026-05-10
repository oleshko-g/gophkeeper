package server

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oleshko-g/gophkeeper/internal/db/pgx"
	"github.com/oleshko-g/gophkeeper/internal/grpc"
	"github.com/oleshko-g/gophkeeper/internal/security"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"google.golang.org/grpc/credentials"
)

// App is the struct to hold a gophkeeper server implementation.
//
// In order to set up the app call:
//
//	[App.I_Configure]
//	[App.II_SetStorage]
//	[App.III_SetSecurity]
//	[App.IV_SetService]
//	[App.V_SetServer]
type App struct {
	grpc struct {
		*grpc.Config
		*grpc.Server
	}
	DB struct {
		*pgx.Config
		*pgx.Conn
	}
	Security struct {
		*security.Config
		securer
	}
	*storage.Storage
	*service.Service
}

type securer interface {
	Certificate() *tls.Certificate
	PrivateKey() *rsa.PrivateKey
	Credentials() *credentials.TransportCredentials
}

// I_Configure loads environment variables and processes them into the [App.config] field.
func (a *App) I_Configure(envFiles ...string) error {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return err
	}

	a.grpc.Config = &grpc.Config{}
	err = envconfig.Process("", a.grpc.Config)
	if err != nil {
		return err
	}

	a.DB.Config = &pgx.Config{}
	err = envconfig.Process("", a.DB.Config)
	if err != nil {
		return err
	}

	a.Service = &service.Service{
		Config: &service.Config{},
	}
	err = envconfig.Process("", a.Service.Config)
	if err != nil {
		return err
	}

	a.Security.Config = &security.Config{}
	err = envconfig.Process("", a.Security.Config)
	if err != nil {
		return err
	}

	return nil
}

// II_SetStorage populates the [App.Storage] field with the given [storage.Keeper] and [storage.Depositor] implementations.
func (a *App) II_SetStorage(depositor storage.Depositor, keeper storage.Keeper) {
	a.Storage = &storage.Storage{
		Depositor: depositor,
		Keeper:    keeper,
	}
}

// III_SetSecurity populates the [App.security.securer] field with the given [security.Config].
func (a *App) III_SetSecurity() error {
	sec, err := security.New(a.Security.Config)
	if err != nil {
		return err
	}

	a.Security.securer = sec
	return nil
}

// IV_SetService populates the [App.Service] field with the given [service.Keeper] and [service.Depositor] implementations.
func (a *App) IV_SetService(depositor service.Depositor, keeper service.Keeper) {
	a.Service = &service.Service{
		Depositor: depositor,
		Keeper:    keeper,
	}
}

// Setup initializes [App.Server] with the [App.Service] implementations. If [App.Service] is nil it panics
func (a *App) V_SetServer(creds credentials.TransportCredentials) error {
	if a.Service == nil || a.Service.Depositor == nil || a.Service.Keeper == nil {
		return errors.New("field Service is nil. SetService() must be called before SetServer")
	}

	srv, err := grpc.New(
		a.Service.Depositor,
		a.Service.Keeper,
		creds,
	)
	if err != nil {
		return err
	}

	a.grpc.Server = srv
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
