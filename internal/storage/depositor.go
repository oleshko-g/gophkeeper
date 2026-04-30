package storage

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
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
	StoreRefreshToken(context.Context, depositor.RefreshToken) error
}
