package storage

import (
	"context"
)

//go:generate moq -rm -out depositor_mock.go . Depositor
type Depositor interface {
	PubKeyer
}

type PubKeyer interface {
	StorePubKey(ctx context.Context, pubKey string) (pubKeyID string, err error)
}
