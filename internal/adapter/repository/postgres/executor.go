package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Executor is the minimal pgx contract required by Postgres repositories.
// It is satisfied by *pgxpool.Pool and pgx.Tx.
type Executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
