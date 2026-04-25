package depositor

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/storage/model"
)

var _ storage.Depositor = (*Depositor)(nil)

type Querier interface {
	InsertRefreshToken(ctx context.Context, arg queries.InsertRefreshTokenParams) (queries.InsertRefreshTokenRow, error)
}

func New(q Querier) *Depositor {
	return &Depositor{
		Querier: q,
	}
}

type Depositor struct {
	Querier
}

func (d *Depositor) StorePubKey(ctx context.Context, pubKey string) (pub_key_id string, err error) {
	if pubKey == "" {
		return "", storage.ErrEmptyInput
	}

	return "not_implemented", nil
}

func (d *Depositor) StoreRefreshToken(ctx context.Context, rt model.RefreshToken) error {
	return nil
}

func (d *Depositor) GetRefreshToken(ctx context.Context, id string) (*model.RefreshToken, error) {
	if id == "" {
		return nil, storage.ErrEmptyInput
	}

	return &model.RefreshToken{}, nil
}
