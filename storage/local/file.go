package local

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

type Storage struct {
	root string
}

const (
	dirPerm  = 0o700
	filePerm = 0o600
)

func New(root string) (*Storage, error) {
	err := os.MkdirAll(root, dirPerm)
	if err != nil {
		return nil, err
	}
	return &Storage{root: root}, nil
}

func NewContent(str string) (string, []byte) {
	content := []byte(str)
	hasher := sha256.New()
	_, err := hasher.Write(content)
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, content
}

func (s *Storage) Write(hash string, reader io.Reader) (bool, error) {
	p := path.Join(s.root, hash)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		if os.IsPermission(err) {
			// file exists
			return true, nil
		}
		return false, err
	}
	defer f.Close()
	hasher := sha256.New()
	rr := io.TeeReader(reader, hasher)
	_, err = io.Copy(f, rr)
	if err != nil {
		os.Remove(p)
		return false, err
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if hash != actual {
		os.Remove(p)
		return false, fmt.Errorf("invalid hash: expected %s, actual %s", hash, actual)
	}
	err = os.Chmod(p, 0o400) // read only
	if err != nil {
		os.Remove(p)
		return false, fmt.Errorf("failed to change file mode: %s", hash)
	}
	return false, nil
}
