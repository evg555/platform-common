package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/evg555/platform-common/pkg/db"
	"github.com/evg555/platform-common/pkg/db/prettier"
)

type key string

// TxKey transaction key for context
const TxKey key = "tx"

type pg struct {
	conn *pgxpool.Pool
}

// NewDB instance of database connection
func NewDB(conn *pgxpool.Pool) db.DB {
	return &pg{
		conn: conn,
	}
}

// ScanOneContext scan query result to single struct
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext scan query result to slice of structs
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext for query
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.conn.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext for query
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.conn.Query(ctx, q.QueryRaw, args...)
}

// QueryRawContext for query
func (p *pg) QueryRawContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.conn.QueryRow(ctx, q.QueryRaw, args...)
}

// Ping database connection
func (p *pg) Ping(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

// Close database connection
func (p pg) Close() {
	p.conn.Close()
}

// BeginTx start transaction
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (db.Committer, error) {
	return p.conn.BeginTx(ctx, txOptions)
}

// MakeContextTx add commiter to context
func MakeContextTx(ctx context.Context, tx db.Committer) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceHolderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
