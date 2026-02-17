package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
)

type ctxKey int

const (
	dbKey                  ctxKey = 0
	defaultMaxIdleConns           = 10
	defaultMaxOpenConns           = 10
	defaultConnMaxLifetime        = 30 * time.Minute
	defaultTimeout                = 5 * time.Second
	defaultPingTimeout            = 3 * time.Second
)

type Config struct {
	Driver          string
	URL             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
}

func (c *Config) maxIdleConns() int {
	if c.MaxIdleConns == 0 {
		return defaultMaxIdleConns
	}
	return c.MaxIdleConns
}

func (c *Config) maxOpenConns() int {
	if c.MaxOpenConns == 0 {
		return defaultMaxOpenConns
	}
	return c.MaxOpenConns
}

func (c *Config) connMaxLifetime() time.Duration {
	if c.ConnMaxLifeTime == 0 {
		return defaultConnMaxLifetime
	}
	return c.ConnMaxLifeTime
}

func NewDB(config *Config) (*sql.DB, error) {
	if config.Driver == "" {
		return nil, fmt.Errorf("database driver is required")
	}
	if config.URL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	d, err := apmsql.Open(config.Driver, config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d.SetMaxIdleConns(config.maxIdleConns())
	d.SetMaxOpenConns(config.maxOpenConns())
	d.SetConnMaxLifetime(config.connMaxLifetime())

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()

	if err = d.PingContext(ctx); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return d, nil
}

func WithTimeout(ctx context.Context, timeout time.Duration, op func(ctx context.Context) error) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return op(ctxWithTimeout)
}

func WithDefaultTimeout(ctx context.Context, op func(ctx context.Context) error) error {
	return WithTimeout(ctx, defaultTimeout, op)
}
