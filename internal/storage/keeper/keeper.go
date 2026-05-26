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
	SelectDepositedSecretIDs(ctx context.Context, depositorPubKeyID uuid.UUID) ([]uuid.UUID, error)
	SelectDepositedSecretData(ctx context.Context, id uuid.UUID) ([]byte, error)
}

func New(q Querier) *Keeper {
	return &Keeper{
		Querier: q,
	}
}

type Keeper struct {
	Querier
}

func (k *Keeper) StoreSecret(ctx context.Context, pkID uuidv7.UUID[uuid.UUID], secretData secret.Data) (*secret.DepositedSecret, error) {
	depositedSecret, err := k.Querier.InsertDepositedSecret(ctx, queries.InsertDepositedSecretParams{
		ID:                uuidv7.New().Value,
		DepositorPubKeyID: pkID.Value,
		EncryptedData:     secretData,
	})
	if err != nil {
		return nil, err
	}

	return secret.FromQueryResult(depositedSecret), nil
}

func (k Keeper) RetrieveSecretIDs(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID]) (uuid.UUIDs, error) {
	return k.Querier.SelectDepositedSecretIDs(ctx, publicKeyID.Value)
}

func (k *Keeper) RetrieveSecret(ctx context.Context, publicKeyID uuidv7.UUID[uuid.UUID]) (secret.Data, error) {
	return k.Querier.SelectDepositedSecretData(ctx, publicKeyID.Value)
}
