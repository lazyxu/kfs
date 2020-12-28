package kfscrypto

import (
	"crypto/cipher"
	"io"
)

type Encoder = cipher.Block

type Compress interface {
	Compress(io.Reader) io.Reader
	Decompress(io.Reader) interface{}
}
