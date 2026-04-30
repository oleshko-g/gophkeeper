package depositor

import (
	"context"

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
	InsertRefreshToken(ctx context.Context, arg queries.InsertRefreshTokenParams) error
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

func (d *Depositor) StoreRefreshToken(ctx context.Context, rt depositor.RefreshToken) error {
	return d.InsertRefreshToken(ctx, queries.InsertRefreshTokenParams{
		Token:             rt.ID.Value,
		DepositorPubKeyID: rt.PubKeyID.Value,
	})
}

func (d *Depositor) GetRefreshToken(ctx context.Context, id string) (*depositor.RefreshToken, error) {
	if id == "" {
		return nil, storage.ErrEmptyInput
	}

	return &depositor.RefreshToken{}, nil
}
