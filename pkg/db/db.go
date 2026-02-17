package db

import (
	"context"
	"database/sql"
	"time"

	"go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
)

type ctxKey int

const (
	dbKey               ctxKey = 0
	defaultMaxIdleConns        = 10
	defaultMaxOpenConns        = 10
	connMaxLifetime            = 30 * time.Minute
	defaultTimeout             = 1 * time.Second
)

var (
	db      *sql.DB
	slaveDB *sql.DB
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

func Init(config *Config) error {
	d, err := NewDB(config)
	if err != nil {
		return err
	}
	db = d
	return nil
}

func InitSlave(config *Config) error {
	d, err := NewDB(config)
	if err != nil {
		return err
	}
	slaveDB = d
	return nil
}

func NewDB(config *Config) (*sql.DB, error) {

	d, err := apmsql.Open(config.Driver, config.URL)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(config.maxIdleConns())
	d.SetMaxOpenConns(config.maxOpenConns())
	d.SetConnMaxLifetime(config.ConnMaxLifeTime)

	return d, err
}

func Close() error {
	return db.Close()
}

func CloseSlave() error {
	return slaveDB.Close()
}

func Get() *sql.DB {
	return db
}

func GetSlave() *sql.DB {
	return slaveDB
}

func WithTimeout(ctx context.Context, timeout time.Duration, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return op(ctxWithTimeout)
}

func WithDefaultTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	return WithTimeout(ctx, defaultTimeout, op)
}
