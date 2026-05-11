package keeper

import (
	"context"

	"github.com/google/uuid"
	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	"github.com/oleshko-g/gophkeeper/internal/model/keeper/secret"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

var _ storage.Keeper = (*Keeper)(nil)

//go:generate moq -rm -out keeper_querier_mock.go . Querier
type Querier interface {
	InsertDepositedSecret(ctx context.Context, arg queries.InsertDepositedSecretParams) (queries.DepositedSecret, error)
}

func New(q Querier) *Keeper {
	return &Keeper{
		Querier: q,
	}
}

type Keeper struct {
	Querier
}

func (k *Keeper) StoreSecret(ctx context.Context, pkID uuidv7.UUID[uuid.UUID], sv secret.Data) (*secret.DepositedSecret, error) {
	depositedSecret, err := k.Querier.InsertDepositedSecret(ctx, queries.InsertDepositedSecretParams{
		ID:                uuidv7.New().Value,
		DepositorPubKeyID: pkID.Value,
		EncryptedData:     sv,
	})
	if err != nil {
		return nil, err
	}

	return secret.FromQueryResult(depositedSecret), nil
}
