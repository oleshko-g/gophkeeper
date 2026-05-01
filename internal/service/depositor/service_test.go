package depositor_test

import (
	"testing"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

func TestAuthorize(t *testing.T) {
	svc := depositor.New(&storage.DepositorMock{}, refreshTokenTTL)

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
	svc := depositor.New(&storage.DepositorMock{}, refreshTokenTTL)

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
