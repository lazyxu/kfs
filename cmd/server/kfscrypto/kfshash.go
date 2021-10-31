package kfscrypto

import (
	"hash"
	"io"
)

type Hash interface {
	// Write (via the embedded io.Writer interface) adds more data to the running hash.
	// It never returns an error.
	io.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Cal(r io.Reader) ([]byte, error)

	// Reset resets the Hash to its initial state.
	Reset()

	// Size returns the number of bytes Sum will return.
	Size() int

	// BlockSize returns the hash's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if all writes
	// are a multiple of the block size.
	BlockSize() int
}

type wrapper struct {
	hash.Hash
}

func (h *wrapper) Cal(r io.Reader) ([]byte, error) {
	//b := new(bytes.Buffer)
	//rr := io.TeeReader(r, b)
	if r != nil {
		_, err := io.Copy(h, r)
		if err != nil {
			return nil, err
		}
	}
	key := h.Sum(nil)
	//fmt.Println("---", key, b.String())
	return key, nil
}

func FromStdHash(h hash.Hash) Hash {
	return &wrapper{h}
}
