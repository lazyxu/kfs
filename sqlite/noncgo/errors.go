package noncgo

import (
	"errors"

	"modernc.org/sqlite"
)

var (
	ErrExpectedDir = errors.New("expected dir")
)

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*sqlite.Error); ok {
		return e.Code() == 1555
	}
	return false
}
