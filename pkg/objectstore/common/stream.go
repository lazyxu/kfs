package common

type ReaderWriter interface {
	Write([]byte) error
	Read(*[]byte) (int, error)
	Close() error
}

type Reader interface {
	Read(*[]byte) (int, error)
	Close() error
}
