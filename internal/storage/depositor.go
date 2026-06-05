package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

type Depositor interface {
	PubKeyer
}

type PubKeyer interface {
	StorePubKey(ctx context.Context, pubKey string) (pubKeyID string, err error)
	RetrievePubKeyByID(ctx context.Context, id uuidv7.UUID[uuid.UUID]) (*depositor.PubKey, error)
}
