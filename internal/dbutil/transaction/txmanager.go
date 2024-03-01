package transaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pintoter/warehouse-api/internal/dbutil"
)

type Manager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) dbutil.TxManager {
	return &Manager{db: db}
}

func (m *Manager) WithTx(ctx context.Context, fn dbutil.Handler) error {
	txOpts := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	return m.transaction(ctx, txOpts, fn)
}

func (m *Manager) transaction(ctx context.Context, txOpts sql.TxOptions, fn dbutil.Handler) (err error) {
	tx, ok := ctx.Value(dbutil.TxKey).(*sql.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, &txOpts)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, dbutil.TxKey, tx)

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		if nil == err {
			_ = tx.Commit()
		}
	}()

	return fn(ctx)
}
