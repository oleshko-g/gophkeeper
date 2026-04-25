package depositor

import (
	"context"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/storage/model"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

var _ service.Depositor = (*Service)(nil)

func New(s storage.Depositor) *Service {
	return &Service{Depositor: s}
}

type Service struct {
	storage.Depositor
}

// Register registers an anonymous public key and returns a refresh token.
// The owner of the refresh token can then [Authorize] apps to [Connect] to [KeeperService]
func (s *Service) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if in == nil {
		return nil, errRegisterRequestIsEmpty
	}

	pubKeyID, err := s.Depositor.StorePubKey(ctx, in.GetPubKey())
	if err != nil {
		return nil, err
	}

	rt := model.RefreshToken{
		PubKeyID:     pubKeyID,
		RefreshToken: uuidv7.NewString(),
	}
	err = s.Depositor.StoreRefreshToken(ctx, rt)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{RefreshToken: &rt.RefreshToken}, nil
}

// Authorize authorizes an app to [Connect] to [KeeperService] and returns an authentication token.
// The owner of the authentication token can then [Connect] to [KeeperService]
func (s *Service) Authorize(ctx context.Context, in *service.AuthorizeRequest) (*service.AuthorizeResponse, error) {
	return nil, nil
}

// Connect creates a new [Session] for an [Authorize]d client app.
// The holder of the session can then make requests to [KeeperService]
func (s *Service) Connect(ctx context.Context, in *service.ConnectRequest) (*service.ConnectResponse, error) {
	return nil, nil
}
