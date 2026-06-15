package depositor_test

import (
	"context"
	"crypto/rsa"
	"testing"

	"crypto/rand"

	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	. "github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func TestAuthorize(t *testing.T) {
	var (
		svc *Service
	)

	t.Run("Setup", func(t *testing.T) {
		privKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}

		svc = New(
			&storage.DepositorMock{
				RetrievePubKeyByIDFunc: func(_ context.Context, id uuidv7.UUID[uuid.UUID]) (*depositor.PubKey, error) {
					return &depositor.PubKey{ID: id},
						nil
				}},
			privKey,
		)

	})
	t.Run("Valid Public Key ID", func(t *testing.T) {

		pubKeyID := uuidv7.New()
		_, err := svc.Authorize(
			nil,
			&pb.AuthorizeRequest{DecryptedId: new(pubKeyID.Value.String())},
		)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Empty Public Key ID", func(t *testing.T) {
		t.Run("Empty Request", func(t *testing.T) {
			_, err := svc.Authorize(
				nil,
				&pb.AuthorizeRequest{DecryptedId: new("")},
			)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})
}
