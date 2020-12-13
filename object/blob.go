package object

import (
	"io"
)

type Blob struct {
	base   *Obj
	Reader io.Reader
}
