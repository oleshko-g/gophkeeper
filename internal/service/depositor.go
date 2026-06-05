package service

import (
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
)

type Depositor interface {
	Namer
	pb.DepositorServiceServer
	ValidateAuthToken(token string) (*depositor.AuthorizedApp, error)
}
