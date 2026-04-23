package depositor

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/db/pgsql/queries"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/storage/model"
)

type DepositorQuerier interface {
	InsertRefreshToken(ctx context.Context, arg queries.InsertRefreshTokenParams) (queries.InsertRefreshTokenRow, error)
}

func New() *Depositor {
	return &Depositor{}
}

var _ storage.Depositor = (*Depositor)(nil)

type Depositor struct {
	DepositorQuerier
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
