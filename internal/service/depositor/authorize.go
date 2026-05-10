package depositor

import (
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/model/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func (s *Service) Authorize(ctx context.Context, in *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	methodName := "Authorize"

	pk, err := s.RetrievePubKeyByID(ctx, uuidv7.FromString(in.GetDecryptedId()))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, s.wrapError(methodName, service.ErrTypeBusinessLogic, err)
		}
		return nil, s.wrapError(methodName, service.ErrTypeStorage, err)
	}

	clientApp := depositor.AuthorizedApp{
		ID: uuidv7.New(),
		PubKey: depositor.PubKey{
			ID: pk.ID,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    s.Name,
				IssuedAt:  jwt.NewNumericDate(clientApp.AuthorizedAt()),
				ExpiresAt: jwt.NewNumericDate(clientApp.AuthorizedAt().Add(24 * time.Hour)),
			},
			AuthorizedApp: clientApp,
		},
	)

	signedToken, err := token.SignedString(s.privKey)

	return &pb.AuthorizeResponse{
		AuthToken: new(signedToken),
	}, nil
}

func (s *Service) ValidateAuthToken(token string) (*depositor.AuthorizedApp, error) {
	methodName := "ValidateAuthToken"

	clms := claims{}
	jwtToken, err := jwt.ParseWithClaims(token, &clms, func(token *jwt.Token) (any, error) {
		return s.privKey.Public(), nil
	})
	if err != nil {
		return nil, s.wrapError(methodName, service.ErrTypeUnauthenticated, err)
	}
	if !jwtToken.Valid {
		return nil, s.wrapError(methodName, service.ErrTypeUnauthenticated, errInvalidToken)
	}

	return &clms.AuthorizedApp, nil
}

type claims struct {
	jwt.RegisteredClaims
	depositor.AuthorizedApp
}
