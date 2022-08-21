package local

import "testing"

func TestNewStorage1(t *testing.T) {
	testNew(t, func(root string) (Storage, error) {
		return NewStorage1(root)
	})
}

func TestErrorHashStorage1(t *testing.T) {
	testErrorHash(t, func(root string) (Storage, error) {
		return NewStorage1(root)
	})
}

func TestWriteTwiceStorage1(t *testing.T) {
	testWriteTwice(t, func(root string) (Storage, error) {
		return NewStorage1(root)
	})
}

func TestConcurrentWriteStorage1(t *testing.T) {
	testConcurrentWrite(t, func(root string) (Storage, error) {
		return NewStorage1(root)
	})
}
