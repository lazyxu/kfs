package kfscrypto

import (
	"crypto/cipher"
	"encoding/gob"
	"io"
)

type Serializable interface {
	Serialize(v interface{}, w io.Writer) error
	Deserialize(v interface{}, r io.Reader) error
}

type Encoder = cipher.Block

type Compress interface {
	Compress(io.Reader) io.Reader
	Decompress(io.Reader) interface{}
}

type GobEncoder struct {
}

func (e *GobEncoder) Serialize(v interface{}, w io.Writer) error {
	return gob.NewEncoder(w).Encode(v)
}

func (e *GobEncoder) Deserialize(v interface{}, r io.Reader) error {
	return gob.NewDecoder(r).Decode(v)
}
