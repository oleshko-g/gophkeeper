package storage

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/storage/model"
)

//go:generate moq -rm -out depositor_mock.go . Depositor
type Depositor interface {
	PubKeyer
	RefreshTokener
}

type PubKeyer interface {
	StorePubKey(ctx context.Context, pubKey string) (pubKeyID string, err error)
}

type RefreshTokener interface {
	StoreRefreshToken(context.Context, model.RefreshToken) error
	GetRefreshToken(ctx context.Context, id string) (*model.RefreshToken, error)
}
