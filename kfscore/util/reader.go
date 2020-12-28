package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func Skip(reader io.Reader, off int64) (int, error) {
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		return int(n), err
	}
	n, err := io.CopyN(ioutil.Discard, reader, off)
	return int(n), err
}

func Size(r io.Reader) (int64, error) {
	switch v := r.(type) {
	case *bytes.Buffer:
		return int64(v.Len()), nil
	case *bytes.Reader:
		return int64(v.Len()), nil
	case *strings.Reader:
		return int64(v.Len()), nil
	case *os.File:
		info, err := v.Stat()
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	default:
		return 0, errors.New("invalid type: " + reflect.TypeOf(r).String())
	}
}
