package server_test

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"time"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	"github.com/oleshko-g/gophkeeper/internal/security"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"google.golang.org/grpc/credentials"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestApp(t0 *testing.T) {
	app := server.App{}
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
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	t0.Run("Configure", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := app.Configure("testdata/.env")
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("Err", func(t *testing.T) {
			err := app.Configure("no env file")
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})

	t0.Run("SetStorage", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}
		app.SetStorage(&storage.DepositorMock{}, storage.KeeperMock{})
	})

	t0.Run("SetService", func(t *testing.T) {
		app := app
		if t0.Failed() {
			t.Skip()
		}
		app.SetService(&service.DepositorMock{}, &service.KeeperMock{})
	})

	t0.Run("SetServer", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}

		t.Run("Success", func(t *testing.T) {
			app := app // make a copy for the test
			app.SetService(&service.DepositorMock{}, &service.KeeperMock{})
		})

		t.Run("Err", func(t *testing.T) {
			err := app.SetServer(creds)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})

	t0.Run("Run", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}

		app := app
		app.Configure("testdata/.env")
		app.SetService(&service.DepositorMock{}, &service.KeeperMock{})

		t.Run("Success", func(t *testing.T) {
			err := app.SetServer(creds)
			if err != nil {
				t0.Error(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			err = app.Run(ctx)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Error(err)
			}
		})

		t.Run("Err", func(t *testing.T) {
			app := app
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := app.Run(ctx)
			if err == nil {
				cancel()
				t.Fatal("expected error, got nil")
			}
		})
	})
}
