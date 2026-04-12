package grpc

import "time"

type Config struct {
	Host       string
	Port       string
	SessionTTL time.Duration
}
