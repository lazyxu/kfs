package local

import "testing"

func TestNewStorage0(t *testing.T) {
	testNew(t, func(root string) (Storage, error) {
		return NewStorage0(root)
	})
}

func TestErrorHashStorage0(t *testing.T) {
	testErrorHash(t, func(root string) (Storage, error) {
		return NewStorage0(root)
	})
}

func TestWriteTwiceStorage0(t *testing.T) {
	testWriteTwice(t, func(root string) (Storage, error) {
		return NewStorage0(root)
	})
}

func TestConcurrentWriteStorage0(t *testing.T) {
	testConcurrentWrite(t, func(root string) (Storage, error) {
		return NewStorage0(root)
	})
}
