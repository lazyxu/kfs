package hashfunc

type HashFunc interface {
	Hash(data []byte) ([]byte, error)
	Size() int
}

const (
	HASH_SHA256 = "sha256"
)

func GetHashFunc(hash string) HashFunc {
	switch hash {
	case HASH_SHA256:
		return &sha256Hash{}
	}
	return nil
}
