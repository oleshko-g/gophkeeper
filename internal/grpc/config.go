package grpc

import "time"

type Config struct {
	GRPCAddr   string        `envconfig:"GRPC_ADDR" default:":8001"`
	SessionTTL time.Duration `envconfig:"SESSION_TTL" default:"15m"`
}
