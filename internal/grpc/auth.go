package grpc

import (
	"context"
	"path"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// authOption checks if ctx has "authorization" metadata.
//   - If ctx has "authorization" authOption authenticates "authorization" first value.
func (s *Server) authOption() grpc.ServerOption {
	return grpc.UnaryInterceptor(
		func(ctx context.Context, req any, srvInfo *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (any, error) {

			if srvInfo.FullMethod == pb.DepositorService_Register_FullMethodName ||
				srvInfo.FullMethod == pb.DepositorService_Authorize_FullMethodName {

				return handler(ctx, req)
			}

			var (
				md metadata.MD
			)

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.New(map[string]string{})
			}

			authVals := md.Get("authorization")
			var authTokenValue string
			if len(authVals) > 0 {
				authTokenValue = authVals[0]
			} else {
				return nil, service.WrapError(&service.Err{
					SvcName: "",
					Method:  path.Base(srvInfo.FullMethod),
					Type:    service.ErrTypeUnauthenticated,
				}, errUnauthenticated)
			}

			authorizedApp, err := s.ValidateAuthToken(authTokenValue)
			if err != nil {
				return nil, err
			}

			md.Set("private_key_id", authorizedApp.PubKey.ID.String)

			return handler(ctx, req)
		})
}
