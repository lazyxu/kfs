package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

type TxOrDb interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type DbImpl interface {
	InsertCommitWithTxOrDb(ctx context.Context, txOrDb TxOrDb, commit *dao.Commit) error
	UpsertBranchWithTxOrDb(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error

	IsUniqueConstraintError(error) bool
	MaxBatchSize() int
}
