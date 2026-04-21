package depositor_test

import (
	"context"
	"testing"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/transform"
)

func TestDepositor(t *testing.T) {
	svc := depositor.New(&storage.DepositorMock{})

	t.Run("Register", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			_, err := svc.Register(
				context.TODO(),
				&pb.RegisterRequest{PubKey: transform.ValueToPtr("pub_key")},
			)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("Err", func(t *testing.T) {
			_, err := svc.Register(
				context.TODO(),
				nil,
			)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	})

	t.Run("success", func(t *testing.T) {

	})
}
