package filesystem

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/pkg/hashfunc"

	"github.com/lazyxu/kfs/pkg/objectstore/common"

	"github.com/sirupsen/logrus"
)

// Store stores data in a filesystem
type Store struct {
	Root   string
	Hash   hashfunc.HashFunc
	Endian binary.ByteOrder
}

func New(root string, hash hashfunc.HashFunc, endian binary.ByteOrder) *Store {
	return &Store{
		Root:   root,
		Hash:   hash,
		Endian: endian,
	}
}

func (s *Store) getPathByHash(hash []byte, createIfNotExist bool) (string, error) {
	if createIfNotExist {
		dir := filepath.Join(s.Root, "objects",
			hex.EncodeToString(hash[0:1]),
			hex.EncodeToString(hash[1:2]),
		)
		err := os.MkdirAll(dir, os.FileMode(0700))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":            err.Error(),
				"hash":             hex.EncodeToString(hash),
				"createIfNotExist": createIfNotExist,
			}).Error("getPathByHash")
			return "", err
		}
		name := hex.EncodeToString(hash[2:])
		return filepath.Join(dir, name), nil
	}
	return filepath.Join(s.Root, "objects",
		hex.EncodeToString(hash[0:1]),
		hex.EncodeToString(hash[1:2]),
		hex.EncodeToString(hash[2:]),
	), nil
}

func (s *Store) GetStream(hash []byte) (r common.Reader, err error) {
	logrus.WithFields(logrus.Fields{
		"hash": hex.EncodeToString(hash),
	}).Debug("GetStream")
	path, err := s.getPathByHash(hash, false)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"hash":  hex.EncodeToString(hash),
			"path":  path,
		}).Error("NewStreamOpen")
		return nil, fmt.Errorf("failed to get object: %s", hex.EncodeToString(hash))
	}
	r = &fileReader{
		file: f,
	}
	return r, nil
}

func (s *Store) PutStream(hash []byte, r common.Reader) error {
	logrus.WithFields(logrus.Fields{
		"hash": hex.EncodeToString(hash),
	}).Debug("PutStream")
	path, _ := s.getPathByHash(hash, true)
	f, err := os.Open(path)
	// if exists, do nothing
	if !os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"path": "path",
		}).Debug("exist")
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	packet := make([]byte, 100*1000)
	for {
		_, readErr := r.Read(&packet)
		if len(packet) == 0 {
			if readErr != io.EOF && readErr != nil {
				logrus.WithFields(logrus.Fields{
					"err": readErr.Error(),
				}).Debug("Read")
				return fmt.Errorf("failed to receive the packet: %s", readErr.Error())
			}
		}
		_, err := f.Write(packet)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": "path",
			}).Debug("exist")
			return readErr
		}
		if readErr == io.EOF {
			break
		}
	}
	return err
}
