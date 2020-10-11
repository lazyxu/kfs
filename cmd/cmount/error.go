package main

import (
	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/lazyxu/kfs/core/e"
	"github.com/sirupsen/logrus"
)

// Translate errors
func translateError(err error) (errc int) {
	defer e.Trace(logrus.Fields{
		"err": err,
	})(func() logrus.Fields {
		return logrus.Fields{
			"errc": errc,
		}
	})
	if err == nil {
		return 0
	}
	switch err {
	case e.OK:
		return 0
	case e.ErrNotExist:
		return -fuse.ENOENT
	case e.ErrExist:
		return -fuse.EEXIST
	case e.ErrPermission:
		return -fuse.EPERM
	case e.ErrClosed:
		return -fuse.EBADF
	case e.ENOTEMPTY:
		return -fuse.ENOTEMPTY
	case e.ESPIPE:
		return -fuse.ESPIPE
	case e.EBADF:
		return -fuse.EBADF
	case e.EROFS:
		return -fuse.EROFS
	case e.ENotImpl:
		return -fuse.ENOSYS
	case e.ErrInvalid:
		return -fuse.EINVAL
	case e.ENotFile:
		return -fuse.ENFILE
	}
	return -fuse.EIO
}
