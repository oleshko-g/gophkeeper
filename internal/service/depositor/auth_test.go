package depositor_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	. "github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func TestAuthorize(t *testing.T) {
	pubKeyID := uuidv7.New()
	svc := New(&storage.DepositorMock{
		RetrievePubKeyByIDFunc: func(ctx context.Context, id uuidv7.UUID[uuid.UUID]) (*depositor.PubKey, error) {
			return &depositor.PubKey{ID: id}, nil
		},
	}, nil)

	t.Run("Success", func(t *testing.T) {
		_, err := svc.Authorize(
			nil,
			&pb.AuthorizeRequest{DecryptedId: new(pubKeyID.Value.String())},
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
