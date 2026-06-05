package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

type Keeper interface {
	StoreSecret(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID], s secret.Data) (*secret.DepositedSecret, error)
	RetrieveSecretIDs(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID]) (uuid.UUIDs, error)
	RetrieveSecret(ctx context.Context, publicKeyID, secretID uuidv7.UUID[uuid.UUID]) (secret.Data, error)
	RemoveSecret(ctx context.Context, publicKeyID, secretID uuidv7.UUID[uuid.UUID]) error
}
