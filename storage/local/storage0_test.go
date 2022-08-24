package local

import "testing"

func TestNewStorage0(t *testing.T) {
	testNew(t, NewStorage0)
}

func TestErrorHashStorage0(t *testing.T) {
	testErrorHash(t, NewStorage0)
}

func TestWriteTwiceStorage0(t *testing.T) {
	testWriteTwice(t, NewStorage0)
}

func TestConcurrentWriteStorage0(t *testing.T) {
	testConcurrentWrite(t, NewStorage0)
}
