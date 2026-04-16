package keeper

import "github.com/oleshko-g/gophkeeper/internal/storage"

func New() *Keeper {
	return &Keeper{}
}

var _ storage.Keeper = (*Keeper)(nil)

type Keeper struct{}
