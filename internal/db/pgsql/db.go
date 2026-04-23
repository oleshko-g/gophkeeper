package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var driverName string = "pgx"

func New(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(driverName, cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db, err = newDB(cfg)
		if err != nil {
			return nil, err
		}
	}

	if err = up(db); err != nil {
		return nil, err
	}
	return db, nil
}

// newDB creates a new database with the name specified in the config and returns a connection to it
func newDB(cfg *Config) (*sql.DB, error) {
	ctx := context.Background()

	defaultDB, err := sql.Open(driverName, cfg.PostgresDefaultDSN)
	defer defaultDB.Close()
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("CREATE DATABASE %s;", cfg.GophkeeperDBName)
	_, err = defaultDB.ExecContext(ctx, q)
	if err != nil {
		return nil, err
	}

	return sql.Open(driverName, cfg.PostgresDSN)
}
