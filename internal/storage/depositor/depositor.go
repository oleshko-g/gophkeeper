package depositor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

var _ storage.Depositor = (*Depositor)(nil)

func New(q Querier) *Depositor {
	return &Depositor{
		Querier: q,
	}
}

type Depositor struct {
	Querier
}

//go:generate moq -rm -out depositor_querier_mock.go . Querier
type Querier interface {
	InsertPubKey(ctx context.Context, arg queries.InsertPubKeyParams) error
	SelectPubKeyByID(ctx context.Context, id uuid.UUID) (queries.DepositorPubKey, error)
}

func (d *Depositor) StorePubKey(ctx context.Context, pubKey string) (pub_key_id string, err error) {
	if pubKey == "" {
		return "", storage.ErrEmptyInput
	}

	id := uuidv7.New()
	err = d.InsertPubKey(ctx, queries.InsertPubKeyParams{
		ID:     id.Value,
		PubKey: pubKey,
	})
	if err != nil {
		return "", err
	}

	return id.Value.String(), nil
}

func (d *Depositor) RetrievePubKeyByID(ctx context.Context, id uuidv7.UUID[uuid.UUID]) (*depositor.PubKey, error) {
	pk, err := d.SelectPubKeyByID(ctx, id.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &depositor.PubKey{
		ID:    uuidv7.FromString(pk.ID.String()),
		Value: pk.PubKey,
	}, nil
}
