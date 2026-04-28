package depositor_test

import (
	"context"
	"errors"
	"testing"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/storage/model"
	"github.com/oleshko-g/gophkeeper/internal/transform"
)

func TestRegister(t *testing.T) {
	svc := depositor.New(&storage.DepositorMock{
		StorePubKeyFunc: func(_ context.Context, _ string) (string, error) {
			return "019dd2b5-0ab9-768b-b1f9-aac25f94d238", nil
		},
		StoreRefreshTokenFunc: func(_ context.Context, _ model.RefreshToken) error {
			return nil
		},
	})
	t.Run("Success", func(t *testing.T) {
		_, err := svc.Register(
			nil,
			&pb.RegisterRequest{PubKey: transform.ValueToPtr("pub_key")},
		)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Empty Request", func(t *testing.T) {
			_, err := svc.Register(
				nil,
				nil,
			)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})

		t.Run("Invalid PubKey", func(t *testing.T) {
			_, err := svc.Register(
				nil,
				&pb.RegisterRequest{PubKey: transform.ValueToPtr("")},
			)
			if err == nil {
				t.Errorf("expected error, got nil")
			}

			e, ok := errors.AsType[*service.Err](err)
			if !ok {
				t.Errorf("expected service.Err, got %v", e)
			}

			if e.SvcName != "Depositor" {
				t.Errorf("expected SvcName depositor, got %s", e.SvcName)
			}

			if e.Method != "Register" {
				t.Errorf("expected method Register, got %s", e.Method)
			}

			pbErr, ok := errors.AsType[pb.RegisterRequestValidationError](e.Err)
			if !ok {
				t.Errorf("expected RegisterRequestValidationError, got %v", pbErr)
			}

			if pbErr.Field() != "PubKey" {
				t.Errorf("expected field pub_key, got %s", pbErr.Field())
			}
		})
	})
}

func TestAuthorize(t *testing.T) {
	svc := depositor.New(&storage.DepositorMock{})

	t.Run("Success", func(t *testing.T) {
		_, err := svc.Authorize(
			nil,
			&pb.AuthorizeRequest{},
		)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Empty Request", func(t *testing.T) {
			_, err := svc.Authorize(
				nil,
				nil,
			)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})
}

func TestConnect(t *testing.T) {
	svc := depositor.New(&storage.DepositorMock{})

	t.Run("Success", func(t *testing.T) {
		_, err := svc.Connect(
			nil,
			&pb.ConnectRequest{},
		)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Empty Request", func(t *testing.T) {
			_, err := svc.Connect(
				nil,
				nil,
			)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})
}
