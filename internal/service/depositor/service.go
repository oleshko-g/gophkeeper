package depositor

import (

	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"

	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"

)

var _ service.Depositor = (*Service)(nil)

func New(s storage.Depositor, refreshTokenTTL time.Duration) *Service {
	return &Service{Depositor: s, refreshTokenTTL: refreshTokenTTL}
}

type Service struct {
	storage.Depositor
	refreshTokenTTL time.Duration
	pb.UnimplementedDepositorServiceServer
}

// Authorize authorizes an app to [Connect] to [KeeperService] and returns an authentication token.
// The owner of the authentication token can then [Connect] to [KeeperService]
// func (s *Service) Authorize(ctx context.Context, in *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
// 	rt, err := s.Depositor.GetRefreshToken(ctx, in.GetRefreshToken())
// 	if err != nil {
// 		return nil, err
// 	}

// 	_ = rt
// 	// if revoked return service.Err

// 	// if expired return service.Err

// 	// var a authToken
// 	// make auth token
// 	// return a

// 	return &pb.AuthorizeResponse{AuthToken: transform.ValueToPtr("")}, nil
// }

// // Connect creates a new [Session] for an [Authorize]d client app.
// // The holder of the session can then make requests to [KeeperService]
// func (s *Service) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.ConnectResponse, error) {
// var a authToken
// 	if err := a.parse(in.GetRefreshToken()); err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }

type authToken struct {
}

func (a *authToken) parse(token string) error {
	return nil
}

func validate() error {
	return nil
}

func authenticate(token string) error {
	return nil
}
