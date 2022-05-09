package local

import (
	"context"
	"errors"
	"fmt"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type Backup struct {
	db *sqlite.SqliteNonCgoDB
}

func New(dbFileName string) (*Backup, error) {
	db, err := sqlite.Open(dbFileName)
	if err != nil {
		return nil, err
	}
	return &Backup{db: db}, nil
}

func (b *Backup) Upload(ctx context.Context, root string) error {
	backupCtx := storage.NewBackupCtx(ctx, root, &uploadVisitor{backup: b})
	ret, err := backupCtx.Scan()
	if err != nil {
		return err
	}
	if dir, ok := ret.(sqlite.Dir); ok {
		status := backupCtx.GetStatus()
		fmt.Printf("%+v\n", status)
		err = b.db.WriteBranch(ctx, sqlite.NewBranch("default", fmt.Sprintf("%+v\n", status), dir))
		if err != nil {
			return err
		}
	} else {
		return errors.New("expected a directory ")
	}
	return nil
}

func (b *Backup) Close() error {
	return b.db.Close()
}
