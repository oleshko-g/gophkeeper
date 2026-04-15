package keeper

import "github.com/oleshko-g/gophkeeper/internal/storage"

func NewKeeper() *Keeper {
	return &Keeper{}
}

var _ storage.Keeper = (*Keeper)(nil)

type Keeper struct{}
