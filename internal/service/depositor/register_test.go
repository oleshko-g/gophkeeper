package depositor_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	. "github.com/oleshko-g/gophkeeper/internal/model/depositor"
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
	"github.com/oleshko-g/gophkeeper/internal/transform"
)

var refreshTokenTTL time.Duration = time.Hour * 24

func TestRegister(t *testing.T) {
	type testCase struct {
		name string
		req  *pb.RegisterRequest
		want *service.Err
	}

	var (
		svc        *depositor.Service
		methodName string
		priv       *rsa.PrivateKey
		err        error
		tests      []testCase
	)

	t.Run("Setup", func(t *testing.T) {
		t.Run("Service", func(t *testing.T) {
			svc = depositor.New(&storage.DepositorMock{
				StorePubKeyFunc:       func(_ context.Context, _ string) (string, error) { return "019dd2b5-0ab9-768b-b1f9-aac25f94d238", nil },
				StoreRefreshTokenFunc: func(_ context.Context, _ RefreshToken) error { return nil }},
				refreshTokenTTL,
			)
		})

		t.Run("Test cases", func(t *testing.T) {
			priv, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatal(err)
			}

			var (
				testName      string
				stringBuilder = &strings.Builder{}
			)

			testName = "Emtpy Request"
			t.Run(testName, func(t *testing.T) {
				tests = append(tests, testCase{name: testName, req: &pb.RegisterRequest{},
					want: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation}})
			})

			testName = "Nil Pub Key"
			t.Run(testName, func(t *testing.T) {
				tests = append(tests, testCase{name: testName, req: &pb.RegisterRequest{PubKey: nil},
					want: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation}})
			})

			testName = "Empty Pub Key"
			t.Run(testName, func(t *testing.T) {
				tests = append(tests, testCase{name: testName, req: &pb.RegisterRequest{PubKey: transform.ValueToPtr("")},
					want: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation}})
			})

			testName = "Empty PKIX Public Key"
			t.Run(testName, func(t *testing.T) {
				err := pem.Encode(stringBuilder, &pem.Block{Type: "PUBLIC KEY", Bytes: nil})
				if err != nil {
					t.Log("encoding")
				}
				tests = append(tests, testCase{name: testName, req: &pb.RegisterRequest{PubKey: transform.ValueToPtr(stringBuilder.String())},
					want: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation}})
				stringBuilder.Reset()
			})

			testName = "Empty PKCS1 Public Key"
			t.Run(testName, func(t *testing.T) {
				err := pem.Encode(stringBuilder, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: nil})
				if err != nil {
					t.Log("encoding")
				}
				tests = append(tests, testCase{name: testName, req: &pb.RegisterRequest{PubKey: transform.ValueToPtr(stringBuilder.String())},
					want: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation}})
				stringBuilder.Reset()
			})

			testName = "Valid PKIX Public Key"
			t.Run(testName, func(t *testing.T) {
				PKIXPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
				err := pem.Encode(stringBuilder, &pem.Block{Type: "PUBLIC KEY", Bytes: PKIXPublicKeyBytes})
				if err != nil {
					t.Log("encoding")
				}
				tests = append(tests, testCase{name: testName,
					req: &pb.RegisterRequest{PubKey: transform.ValueToPtr(stringBuilder.String())}, want: nil})
				stringBuilder.Reset()
			})

			testName = "Valid PKCS1 Public Key"
			t.Run(testName, func(t *testing.T) {
				err = pem.Encode(stringBuilder, &pem.Block{
					Type:  "RSA PUBLIC KEY",
					Bytes: x509.MarshalPKCS1PublicKey(&priv.PublicKey),
				})
				if err != nil {
					t.Log("encoding")
				}

				tests = append(tests, testCase{name: testName,
					req: &pb.RegisterRequest{PubKey: transform.ValueToPtr(stringBuilder.String())}, want: nil})
				stringBuilder.Reset()
			})
		})
	})

	t.Run("Tests", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := svc.Register(nil, test.req)
				svcErr, ok := errors.AsType[*service.Err](err)
				if !ok {
					t.Errorf("expected error %v, got %v", test.want, err)
				}

				if test.want.Type != svcErr.Type {
					t.Errorf("expected error type %s, got %s", test.want, err)
				}
			})
		}
	})
}
