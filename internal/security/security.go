package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"

	"google.golang.org/grpc/credentials"
)

func New(cfg *Config) (*Security, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		err = writeX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}
		cert, err = tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}
	}
	s := &Security{
		cert: &cert,
	}

	priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("cert is not an RSA private key")
	}
	s.priv = priv

	creds := credentials.NewServerTLSFromCert(s.cert)
	if err != nil {
		return nil, err
	}
	s.creds = &creds

	return s, nil
}

type Security struct {
	cert  *tls.Certificate
	priv  *rsa.PrivateKey
	creds *credentials.TransportCredentials
}

func (s *Security) Certificate() *tls.Certificate {
	return s.cert
}

func (s *Security) PrivateKey() *rsa.PrivateKey {
	return s.priv
}

func (s *Security) Credentials() *credentials.TransportCredentials {
	return s.creds
}

// createCertificate creates [x509.Certificate] and writes it in the [certFile]
func writeX509KeyPair(certFile, keyFile string) error {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit) //revive:disable-line [crypto/rand.Reader] fills up the buffer and never returns an error

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"gophkeeper"},
			CommonName:   "localhost",
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.ParseIP("::1")},
		DNSNames:              []string{"localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Self-sign: template == parent
	certData, err := x509.CreateCertificate(rand.Reader,
		template,
		template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certData},
	)

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey)},
	)

	if err = os.WriteFile(certFile, certPEM, 0o600); err != nil {
		return err
	}
	if err = os.WriteFile(keyFile, keyPEM, 0o600); err != nil {
		return err
	}

	return nil
}
