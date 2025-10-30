package repository

import (
	"context"
	"database/sql"
)

type TxProvider interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(ctx context.Context, tx *sql.Tx) error
	RollbackTx(ctx context.Context, tx *sql.Tx) error
}
