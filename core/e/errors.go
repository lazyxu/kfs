// Cross platform errors

package e

import (
	"fmt"
	"os"
)

// Error describes low level errors in a cross platform way.
type Error byte

// NB if changing errors translateError in cmd/mount/fs.go, cmd/cmount/fs.go

// Low level errors
const (
	OK Error = iota
	ENotEmpty
	ESPIPE
	EBADF
	EROFS
	ENotImpl
	ENotFile
	ENotDir
	EInvalidType
	EWriteObject
)

// Errors which have exact counterparts in os
var (
	ErrNotExist   = os.ErrNotExist
	ErrExist      = os.ErrExist
	ErrPermission = os.ErrPermission
	ErrInvalid    = os.ErrInvalid
	ErrClosed     = os.ErrClosed
)

var errorNames = []string{
	OK:           "Success",
	ENotEmpty:    "Directory not empty",
	ESPIPE:       "Illegal seek",
	EBADF:        "Bad file descriptor",
	EROFS:        "Read only file system",
	ENotImpl:     "Function not implemented",
	ENotFile:     "Not a file",
	ENotDir:      "Not a directory",
	EInvalidType: "Invalid object type",
	EWriteObject: "Failed to write object",
}

// Error renders the error as a string
func (e Error) Error() string {
	if int(e) >= len(errorNames) {
		return fmt.Sprintf("Low level error %d", e)
	}
	return errorNames[e]
}
