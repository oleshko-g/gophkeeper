package depositor_test

import (
	"context"
	"testing"

	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/storage/depositor"
)

func TestInsertPubKey(t *testing.T) {
	querier := &depositor.QuerierMock{
		InsertPubKeyFunc: func(_ context.Context, _ queries.InsertPubKeyParams) error {
			return nil
		},
		InsertRefreshTokenFunc: func(_ context.Context, _ queries.InsertRefreshTokenParams) error {
			return nil
		},
	}
	storage := depositor.New(querier)

	t.Run("Success", func(t *testing.T) {
		_, err := storage.StorePubKey(nil, "pub_key")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Err", func(t *testing.T) {
		_, err := storage.StorePubKey(nil, "")
		if err == nil {
			t.Error("expected error, got nil", err)
		}
	})
}
