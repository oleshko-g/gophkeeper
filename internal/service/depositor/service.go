package clientapp

import (
	"context"

	"github.com/oleshko-g/gophkeeper/internal/service"
)

var _ service.Depositor = (*Service)(nil)

type Service struct {
}

// Register registers an anonymous public key and returns a refresh token.
// The owner of the refresh token can then [Authorize] apps to [Connect] to [KeeperService]
func (s *Service) Register(ctx context.Context, in *service.RegisterRequest) (*service.RegisterResponse, error) {
	return nil, nil
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
