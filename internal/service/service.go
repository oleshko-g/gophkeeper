package service

import (
	"time"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
)

type Service struct {
	Depositor
	Keeper
	*Config
}

type Config struct {
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL_HRS" default:"24h" `
}

//go:generate moq -rm -out depositor_mock.go . Depositor
type Depositor = pb.DepositorServiceServer

//go:generate moq -rm -out keeper_mock.go . Keeper
type Keeper = pb.KeeperServiceServer
