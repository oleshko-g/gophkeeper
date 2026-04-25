package pgx

import (
	"database/sql"
	"embed"
	_ "embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func up(db *sql.DB) error {
	if err := goose.SetDialect(driverName); err != nil {
		return err
	}

	var dir string
	goose.SetBaseFS(migrations)
	dir = "migrations"

	err := goose.Up(db, dir)
	if err != nil {
		return err
	}

	return nil
}
