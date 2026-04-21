package server_test

import (
	"context"
	"errors"
	"testing"
	"time"
	_ "time"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestApp(t0 *testing.T) {
	app := server.App{}

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
		app.SetStorage(storage.KeeperMock{}, &storage.DepositorMock{})
	})

	t0.Run("SetService", func(t *testing.T) {
		app := app
		if t0.Failed() {
			t.Skip()
		}
		app.SetService(&service.KeeperMock{}, &service.DepositorMock{})
	})

	t0.Run("SetServer", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}

		t.Run("Success", func(t *testing.T) {
			app := app // make a copy for the test
			app.SetService(&service.KeeperMock{}, &service.DepositorMock{})
		})

		t.Run("Err", func(t *testing.T) {
			err := app.SetServer()
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
		app.SetService(&service.KeeperMock{}, &service.DepositorMock{})

		t.Run("Success", func(t *testing.T) {
			err := app.SetServer()
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

			err := app.Run(ctx)
			if err == nil {
				cancel()
				t.Fatal("expected error, got nil")
			}
		})
	})
}
