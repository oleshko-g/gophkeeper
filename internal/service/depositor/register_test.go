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
	"github.com/oleshko-g/gophkeeper/internal/service"
	"github.com/oleshko-g/gophkeeper/internal/service/depositor"
	"github.com/oleshko-g/gophkeeper/internal/storage"
)

var refreshTokenTTL time.Duration = time.Hour * 24

func TestRegister(t *testing.T) {
	type testCase struct {
		req     *pb.RegisterRequest
		wantErr *service.Err
	}

	var (
		svc                         *depositor.Service
		methodName                  string = "Register"
		privServerKey, privCientKey *rsa.PrivateKey
		err                         error
		tests                       = make(map[string]testCase)
	)

	t.Run("Setup", func(t *testing.T) {
		privServerKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}

		privCientKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		t.Run("Service", func(t *testing.T) {
			svc = depositor.New(
				&storage.DepositorMock{
					StorePubKeyFunc: func(_ context.Context, _ string) (string, error) { return "019dd2b5-0ab9-768b-b1f9-aac25f94d238", nil },
				},
				privServerKey,
			)
		})

		t.Run("Test cases", func(t *testing.T) {
			privCientKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatal(err)
			}

			var (
				testName      string
				stringBuilder = &strings.Builder{}
			)

			testName = "Emtpy Request"
			t.Run(testName, func(t *testing.T) {
				tests[testName] = testCase{
					req:     nil,
					wantErr: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation},
				}
			})

			testName = "Nil Pub Key"
			t.Run(testName, func(t *testing.T) {
				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: nil},
					wantErr: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation},
				}
			})

			testName = "Empty Pub Key"
			t.Run(testName, func(t *testing.T) {
				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: new("")},
					wantErr: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation},
				}
			})

			testName = "Empty PKIX Public Key"
			t.Run(testName, func(t *testing.T) {
				err := pem.Encode(stringBuilder, &pem.Block{Type: "PUBLIC KEY", Bytes: nil})
				if err != nil {
					t.Error(err)
				}

				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: new(stringBuilder.String())},
					wantErr: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation},
				}

				stringBuilder.Reset()
			})

			testName = "Empty PKCS1 Public Key"
			t.Run(testName, func(t *testing.T) {
				err := pem.Encode(stringBuilder, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: nil})
				if err != nil {
					t.Error(err)
				}

				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: new(stringBuilder.String())},
					wantErr: &service.Err{SvcName: svc.Name, Method: methodName, Type: service.ErrTypeRequestValidation},
				}

				stringBuilder.Reset()
			})

			testName = "Valid PKIX Public Key"
			t.Run(testName, func(t *testing.T) {
				PKIXPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privCientKey.PublicKey)
				err := pem.Encode(stringBuilder, &pem.Block{Type: "PUBLIC KEY", Bytes: PKIXPublicKeyBytes})
				if err != nil {
					t.Error(err)
				}

				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: new(stringBuilder.String())},
					wantErr: nil,
				}

				stringBuilder.Reset()
			})

			testName = "Valid PKCS1 Public Key"
			t.Run(testName, func(t *testing.T) {
				err = pem.Encode(stringBuilder, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&privCientKey.PublicKey)})
				if err != nil {
					t.Error(err)
				}

				tests[testName] = testCase{
					req:     &pb.RegisterRequest{PubKey: new(stringBuilder.String())},
					wantErr: nil,
				}

				stringBuilder.Reset()
			})
		})
	})

	t.Run("Tests", func(t *testing.T) {
		for testName, test := range tests {
			t.Run(testName, func(t *testing.T) {
				_, err := svc.Register(nil, test.req)
				if test.wantErr != nil {
					svcErr, ok := errors.AsType[*service.Err](err)
					if !ok {
						t.Fatalf("expected: %T, got: %T", test.wantErr, err)
					}

					if test.wantErr.Type != svcErr.Type {
						t.Errorf("expected: %q, got: %q", test.wantErr.Type, svcErr.Type)
					}
				}
			})
		}
	})
}
