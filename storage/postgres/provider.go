package postgres

import (
	"fmt"

	"github.com/Fox1N69/iq-testtask/internal/config"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(config *config.Config) (*pgxpool.Pool, error) {
	client := NewPSQLClient()
	err := client.Connect(config.Psql.User, config.Psql.Password, config.Psql.Host, config.Psql.Port, config.Psql.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return client.DB, nil
}

var ProviderSet = wire.NewSet(NewPostgresDB)
