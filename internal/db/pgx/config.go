package pgx

type Config struct {
	GophkeeperDBName   string `envconfig:"GOPHKEEPER_DB_NAME" default:"gophkeeper"`
	PostgresDefaultDSN string `envconfig:"POSTGRES_DEFAULT_DSN" default:"postgres://postgres:gophkeeper@localhost:5432/postgres?sslmode=disable"`
	PostgresDSN        string `envconfig:"POSTGRES_DSN" default:""`
}
