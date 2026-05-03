package client

type Config struct {
	GophKeeperURL string `envconfig:"GOPHKEEPER_URL" default:":8001"`
}
