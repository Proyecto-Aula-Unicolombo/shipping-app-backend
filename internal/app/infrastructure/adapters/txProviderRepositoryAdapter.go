package adapters

import (
	"context"
	"database/sql"
)

type SQLTxProvider struct {
	db *sql.DB
}

func NewSQLTxProvider(db *sql.DB) *SQLTxProvider {
	return &SQLTxProvider{db: db}
}

func (p *SQLTxProvider) BeginTx(ctx context.Context) (*sql.Tx, error) {

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (p *SQLTxProvider) CommitTx(ctx context.Context, tx *sql.Tx) error {
	if tx == nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *SQLTxProvider) RollbackTx(ctx context.Context, tx *sql.Tx) error {
	if tx == nil {
		return nil
	}

	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		return err
	}

	return nil
}
