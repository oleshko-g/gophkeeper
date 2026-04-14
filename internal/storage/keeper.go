package storage

import "github.com/oleshko-g/gophkeeper/internal/service/keeper"

func NewKeeper() *Keeper {
	return &Keeper{}
}

var _ keeper.Storage = (*keeper.Storage)(nil)

type Keeper struct{}
