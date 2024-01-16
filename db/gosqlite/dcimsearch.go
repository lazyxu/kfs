package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) SearchDCIM(ctx context.Context, typeList []string, suffixList []string) (list []dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.SearchDCIM(ctx, conn, typeList, suffixList)
}
