package dbBase

import "context"

func UpsertLivePhoto(ctx context.Context, txOrDb TxOrDb, movHash string, heicHash string, jpgHash string) (err error) {
	if heicHash != "" && jpgHash != "" {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			heicHash,
			jpgHash
		) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			heicHash=?,
			jpgHash=?;
		`, movHash, heicHash, jpgHash, movHash, heicHash, jpgHash)
	} else if heicHash != "" {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			heicHash
		) VALUES (?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			heicHash=?;
		`, movHash, heicHash, movHash, heicHash)
	} else {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			jpgHash
		) VALUES (?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			jpgHash=?;
		`, movHash, jpgHash, movHash, jpgHash)
	}
	if err != nil {
		return err
	}
	return err
}
