package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func newContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, dbKey, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(dbKey).(*sqlx.Tx)
	if !ok {
		panic("No DB transaction found in context")
	}
	return tx
}

func Transact(ctx context.Context, dbx *sqlx.DB, opts *sql.TxOptions, txFunc func(context.Context) error) (err error) {
	tx, err := dbx.BeginTxx(ctx, opts)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = errors.WithStack(p)
			default:
				err = errors.Errorf("%s", p)
			}
		}
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errors.WithStack(e)
			}
			return
		}
		err = errors.WithStack(tx.Commit())
	}()

	ctxWithTx := newContext(ctx, tx)
	err = WithDefaultTimeout(ctxWithTx, txFunc)
	return err
}

func WithTimeout(ctx context.Context, timeout time.Duration, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return op(ctxWithTimeout)
}

func WithDefaultTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	return WithTimeout(ctx, defaultTimeout, op)
}
