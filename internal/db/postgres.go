package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config - configuration for DataBase.
type Config struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     string `envconfig:"PORT" default:"5432"`
	Name     string `envconfig:"NAME" default:"demo_db"`
	User     string `envconfig:"USER" default:"demo_user"`
	Password string `envconfig:"PASSWORD" default:"demo_password"`
}

// getPsqlDsn generates a PostgreSQL connection string
// based on the provided database configuration.
//
// postgresql://<user>:<password>@<host>:<port>/<name>?sslmode=disable
//
// or
//
// host=<host> port=<port> user=<user> password=<password> dbname=<name> sslmode=disable
func getPsqlDsn(cfg Config) string {
	// dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
	// 	cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	return dsn
}

// NewPostgresPool creates a configuration database pool.
func NewPostgresPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(getPsqlDsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	poolCfg.MinConns = 1                                 // the minimum number of connections
	poolCfg.MaxConnIdleTime = time.Minute * 15           // maximum compound downtime
	poolCfg.MaxConnLifetime = time.Minute * 30           // maximum connection lifetime
	poolCfg.HealthCheckPeriod = time.Minute * 3          // the period of checking the health of the compounds
	poolCfg.ConnConfig.ConnectTimeout = time.Second * 10 // timeout connection

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Connected to DataBase (using pool) successfully!")
	return pool, nil
}
