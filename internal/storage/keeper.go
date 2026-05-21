package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

//go:generate moq -rm -out keeper_mock.go . Keeper
type Keeper interface {
	StoreSecret(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID], s secret.Data) (*secret.DepositedSecret, error)
	RetrieveSecrets(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID]) ([]secret.DepositedSecret, error)
}
