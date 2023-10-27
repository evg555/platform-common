package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler function
type Handler func(ctx context.Context) error

// Client interface for database client
type Client interface {
	DB() DB
	Close() error
}

// Query for query builder
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor interface for transactions
type Transactor interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (Committer, error)
}

// Committer interface for transactions
type Committer interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// TxManager interface for TxManager
type TxManager interface {
	ReadComitted(ctx context.Context, f Handler) error
}

// SQLExecutor interface for execution sql queries
type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

// NamedExecutor interface for scan sql results to structs
type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecutor interface for execution sql queries
type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRawContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for ping database connection
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB common interface for database
type DB interface {
	SQLExecutor
	Pinger
	Transactor
	Close()
}
