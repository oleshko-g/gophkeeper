package storage

import "github.com/oleshko-g/gophkeeper/internal/service/depositor"


func NewDepositor() *Depositor {
	return &Depositor{}
}

var _ depositor.Storage = &Depositor{}

type Depositor struct {

}
