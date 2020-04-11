package common

type Block interface {
	ReadObject(hash []byte) ([]byte, error)
	Read() ([]byte, error)
	Close() error
}
