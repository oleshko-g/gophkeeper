package depositor

import "github.com/oleshko-g/gophkeeper/internal/storage"

func NewDepositor() *Depositor {
	return &Depositor{}
}

var _ storage.Depositor = (*Depositor)(nil)

type Depositor struct {
}
