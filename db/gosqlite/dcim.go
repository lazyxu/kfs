package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) ListDCIMDriver(ctx context.Context) (drivers []dao.DCIMDriver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDCIMDriver(ctx, conn)
}
