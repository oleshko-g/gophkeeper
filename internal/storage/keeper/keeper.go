package keeper

import "github.com/oleshko-g/gophkeeper/internal/storage"

var _ storage.Keeper = (*Keeper)(nil)

//go:generate moq -rm -out keeper_querier_mock.go . Querier
type Querier interface{}

func New(q Querier) *Keeper {
	return &Keeper{
		Querier: q,
	}
}

type Keeper struct {
	Querier Querier
}
