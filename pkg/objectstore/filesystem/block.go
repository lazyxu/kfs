package filesystem

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/lazyxu/kfs/pkg/objectstore/common"
)

func (s *Store) getPath(hash []byte) string {
	return filepath.Join(s.Root, "objects",
		hex.EncodeToString(hash[0:1]),
		hex.EncodeToString(hash[1:2]),
		hex.EncodeToString(hash[2:]),
	)
}

func (s *Store) ReadObject(hash []byte) ([]byte, error) {
	actual := len(hash)
	expected := s.Hash.Size()
	if actual != expected {
		return nil, fmt.Errorf("invalid hash size, expected %d, actual %d", expected, actual)
	}
	filename := s.getPath(hash)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteObject writes object to file system, if object already exists, it returns true and error.
func (s *Store) WriteObject(typ uint32, hash []byte, data []byte, info []byte) (bool, error) {
	actual := len(hash)
	expected := s.Hash.Size()
	if actual != expected {
		return false, fmt.Errorf("invalid hash size, expected %d, actual %d", expected, actual)
	}
	filename := s.getPath(hash)

	err := os.MkdirAll(path.Dir(filename), os.FileMode(0755))
	if err != nil {
		return false, err
	}
	object := common.NewHeader(typ, uint64(len(data)), s.Endian)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.FileMode(0444))
	if os.IsPermission(err) {
		return true, err
	} else if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer([]byte{})
	err = binary.Write(buffer, s.Endian, *object)
	if err != nil {
		return false, fmt.Errorf("failed to write to buffer(%s): %s", filename, err.Error())
	}
	buffer.Write(data)
	buffer.Write(info)
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return false, fmt.Errorf("failed to write file(%s): %s", filename, err.Error())
	}
	f.Close()
	return true, nil
}
