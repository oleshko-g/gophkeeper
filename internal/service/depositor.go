package service

import (
	"context"
)

// Depositor is the interface to:
//   - [Register] public keys of the Keeper users
//   - [Authorize] their client apps
//   - and [Connect] through the authorized apps to [KeeperService]
//
//go:generate moq -out depositor_mock.go . Depositor
type Depositor interface {
	// Register registers an anonymous public key and returns a refresh token.
	// The owner of the refresh token can then [Authorize] apps to [Connect] to [KeeperService]
	Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error)
	// Authorize authorizes an app to [Connect] to [KeeperService] and returns an authentication token.
	// The owner of the authentication token can then [Connect] to [KeeperService]
	Authorize(ctx context.Context, in *AuthorizeRequest) (*AuthorizeResponse, error)
	// Connect creates a new [Session] for an [Authorize]d client app.
	// The holder of the session can then make requests to [KeeperService]
	Connect(ctx context.Context, in *ConnectRequest) (*ConnectResponse, error)
}

type RegisterRequest struct {
}

type RegisterResponse struct {
}

type AuthorizeRequest struct {
}

type AuthorizeResponse struct {
}

type ConnectRequest struct {
}

type ConnectResponse struct {
}
