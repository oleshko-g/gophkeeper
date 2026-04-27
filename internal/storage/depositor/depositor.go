package depositor

import (
	"context"
	"time"

	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/storage/model"
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
		ID:     id,
		PubKey: pubKey,
	})
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (d *Depositor) StoreRefreshToken(ctx context.Context, rt model.RefreshToken) error {
	return d.InsertRefreshToken(ctx, queries.InsertRefreshTokenParams{
		ID:                uuidv7.New(),
		DepositorPubKeyID: uuidv7.FromString(rt.PubKeyID),
		Token:             rt.RefreshToken,
		IssuedAt:          time.Now().UTC(),
	})
}

func (d *Depositor) GetRefreshToken(ctx context.Context, id string) (*model.RefreshToken, error) {
	if id == "" {
		return nil, storage.ErrEmptyInput
	}

	return &model.RefreshToken{}, nil
}
