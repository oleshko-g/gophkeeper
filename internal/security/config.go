package security

// Config holds the configuration for the security package.
type Config struct {
	CertFile string `envconfig:"GOPH_KEEPER_CERT_FILE" default:"cert.pem"`
	KeyFile  string `envconfig:"GOPH_KEEPER_KEY_FILE" default:"key.pem"`
}
