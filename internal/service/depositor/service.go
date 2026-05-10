package depositor

import (
	"crypto/rsa"

	pb "github.com/oleshko-g/gophkeeper/api/v1"

	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

var _ service.Depositor = (*Service)(nil)

func New(s storage.Depositor, priv *rsa.PrivateKey) *Service {
	return &Service{Name: "Depositor", Depositor: s, privKey: priv}
}

type Service struct {
	Name string
	storage.Depositor
	privKey *rsa.PrivateKey
	pb.UnimplementedDepositorServiceServer
}

func (s *Service) wrapError(methodName string, t service.ErrType, err error) error {
	return service.WrapError(&service.Err{SvcName: s.Name, Method: methodName, Type: t}, err)
}
