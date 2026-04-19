package server_test

import (
	"context"
	"testing"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestApp(t0 *testing.T) {
	app := server.App{}

	t0.Run("Configure", func(t1 *testing.T) {
		t1.Run("Success", func(t *testing.T) {
			err := app.Configure("testdata/.env")
			if err != nil {
				t.Fatal(err)
			}
		})

		t1.Run("Err", func(t *testing.T) {
			err := app.Configure("testdata/.env.notfound")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})

	t0.Run("SetStorage", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}
		app.SetStorage(storage.KeeperMock{}, storage.DepositorMock{})
	})

	t0.Run("SetService", func(t *testing.T) {
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
				t.Fatal("expected error, got nil")
			}
		})
	})

	t0.Run("Run", func(t *testing.T) {
		if t0.Failed() {
			t.Skip()
		}
		app := app // make a copy for the test
		app.Configure("testdata/.env")
		app.SetService(&service.KeeperMock{}, &service.DepositorMock{})
		err := app.SetServer()

		t.Run("Success", func(t *testing.T) {
			ctx := context.Background()

			err = app.Run(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Run("Err", func(t *testing.T) {
			ctx := context.Background()
			err := app.Run(ctx)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
}
