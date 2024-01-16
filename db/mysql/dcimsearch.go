package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) SearchDCIM(ctx context.Context, typeList []string, suffixList []string) (list []dao.Metadata, err error) {
	return dbBase.SearchDCIM(ctx, db.db, typeList, suffixList)
}
