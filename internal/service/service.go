package service

import (
	"time"
)

type Service struct {
	Depositor
	Keeper
	*Config
}

type Config struct {
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL_HRS" default:"24h" `
}

type Namer interface {
	Name() string
}
