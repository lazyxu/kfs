package local

import "testing"

func TestNewStorage1(t *testing.T) {
	testNew(t, NewStorage1)
}

func TestErrorHashStorage1(t *testing.T) {
	testErrorHash(t, NewStorage1)
}

func TestWriteTwiceStorage1(t *testing.T) {
	testWriteTwice(t, NewStorage1)
}

func TestConcurrentWriteStorage1(t *testing.T) {
	testConcurrentWrite(t, NewStorage1)
}
