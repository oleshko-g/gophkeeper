package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var driverName string = "pgx"

func New(cfg *Config) (*pgx.Conn, error) {
	pgxConn, err := pgx.Connect(context.Background(), cfg.PostgresDSN)
	if err != nil {
		pgxConn, err = newConn(cfg)
		if err != nil {
			return nil, err
		}
	}

	sqlDB := stdlib.OpenDB(*pgxConn.Config())
	defer sqlDB.Close()

	if err = up(sqlDB); err != nil {
		return nil, err
	}
	return pgxConn, nil
}

// newConn creates a new database with the name specified in the config and returns a connection to it
func newConn(cfg *Config) (*pgx.Conn, error) {
	ctx := context.Background()

	defaultDB, err := pgx.Connect(ctx, cfg.PostgresDefaultDSN)
	defer defaultDB.Close(ctx)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("CREATE DATABASE %s;", cfg.GophkeeperDBName)
	_, err = defaultDB.Exec(ctx, q)
	if err != nil {
		return nil, err
	}

	return pgx.Connect(ctx, cfg.PostgresDSN)
}
