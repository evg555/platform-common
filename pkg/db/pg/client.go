package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/evg555/platform-common/pkg/db"
)

type pgClient struct {
	MasterDB db.DB
}

// New instance of client
func New(ctx context.Context, dsn string) (db.Client, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{MasterDB: NewDB(conn)}, nil
}

// DB return database connection
func (p *pgClient) DB() db.DB {
	return p.MasterDB
}

// Close database connection
func (p *pgClient) Close() error {
	if p.MasterDB != nil {
		p.MasterDB.Close()
	}

	return nil
}
