package mysql

import (
	"context"
	"database/sql"
)

// Executor is the minimal sql contract required by MySQL repositories.
// It is satisfied by *sql.DB and *sql.Tx.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
