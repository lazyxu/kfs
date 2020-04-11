package hashfunc

import "crypto/sha256"

type sha256Hash struct {
}

func (*sha256Hash) Hash(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (*sha256Hash) Size() int {
	return sha256.Size
}
