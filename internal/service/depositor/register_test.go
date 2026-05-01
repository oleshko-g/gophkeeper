package depositor_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	svc := depositor.New(&storage.DepositorMock{
		StorePubKeyFunc: func(_ context.Context, _ string) (string, error) {
			return "019dd2b5-0ab9-768b-b1f9-aac25f94d238", nil
		},
		StoreRefreshTokenFunc: func(_ context.Context, _ RefreshToken) error {
			return nil
		},
	}, refreshTokenTTL)
	t.Run("Success", func(t *testing.T) {
		t.Run("PKIX Public Key", func(t *testing.T) {

			PKIXPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
			pemPKIXPublicKey := pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: PKIXPublicKeyBytes,
			})

			_, err := svc.Register(
				nil,
				&pb.RegisterRequest{PubKey: transform.ValueToPtr(string(pemPKIXPublicKey))},
			)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("PKCS1 Public Key", func(t *testing.T) {
			pemPKCS1PublicKey := pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&priv.PublicKey),
			})
			_, err := svc.Register(
				nil,
				&pb.RegisterRequest{PubKey: transform.ValueToPtr(string(pemPKCS1PublicKey))},
			)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Empty Request", func(t *testing.T) {
			_, err := svc.Register(
				nil,
				nil,
			)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})

		t.Run("Invalid PubKey", func(t *testing.T) {
			_, err := svc.Register(
				nil,
				&pb.RegisterRequest{PubKey: transform.ValueToPtr("")},
			)
			if err == nil {
				t.Errorf("expected error, got nil")
			}

			e, ok := errors.AsType[*service.Err](err)
			if !ok {
				t.Errorf("expected service.Err, got %v", e)
			}

			if e.SvcName != "Depositor" {
				t.Errorf("expected SvcName depositor, got %s", e.SvcName)
			}

			if e.Method != "Register" {
				t.Errorf("expected method Register, got %s", e.Method)
			}

			pbErr, ok := errors.AsType[pb.RegisterRequestValidationError](e.Err)
			if !ok {
				t.Errorf("expected RegisterRequestValidationError, got %v", pbErr)
			}

			if pbErr.Field() != "PubKey" {
				t.Errorf("expected field pub_key, got %s", pbErr.Field())
			}
		})
	})
}
