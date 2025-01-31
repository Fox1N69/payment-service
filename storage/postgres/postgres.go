package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Fox1N69/logger-setup"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PSQLClient struct {
	DB  *pgxpool.Pool
	log logger.Logger
}

func NewPSQLClient() *PSQLClient {
	return &PSQLClient{
		log: logger.GetLogger(),
	}
}

func (s *PSQLClient) Connect(user, password, host, port, dbname string) error {
	const op = "storage.postgres.Connect()"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("%s: failed to parse DSN: %w", op, err)
	}

	config.MaxConns = 750
	config.MinConns = 10

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	s.DB = dbpool

	return nil
}

func (s *PSQLClient) Close() {
	if s.DB != nil {
		s.log.Info("Closing the database connection...")
		s.DB.Close()
		s.log.Info("Database connection closed successfully")
	} else {
		s.log.Warn("Attempted to close a nil database connection")
	}
}
