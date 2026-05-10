package grpc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big" // Need this for Serial Number
	"net"
	"os"
	"time"
)

const (
	certFile, keyFile string = "cert.pem", "key.pem"
)

func newTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = writeX509KeyPair()
		if err != nil {
			return nil
		}

		cert, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil
		}
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

// createCertificate creates [x509.Certificate] and writes it in the [certFile]
func writeX509KeyPair() error {
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
