package kfshash

import (
	"encoding/hex"
	"hash"
	"io"
)

type Hash interface {
	// Write (via the embedded io.Writer interface) adds more data to the running hash.
	// It never returns an error.
	io.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Cal(r io.Reader) (string, error)

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

func (h *wrapper) Cal(r io.Reader) (string, error) {
	if _, err := io.Copy(h, r); err != nil {
		return "", nil
	}
	key := hex.EncodeToString(h.Sum(nil))
	return key, nil
}

func FromStdHash(h hash.Hash) Hash {
	return &wrapper{h}
}
