package node

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/object"

	"github.com/sirupsen/logrus"
)

type File struct {
	ItemBase
}

func NewFile(s storage.Storage, obj *object.Obj, metadata *object.Metadata, parent *Dir) *File {
	return &File{
		ItemBase: ItemBase{
			storage:  s,
			obj:      obj,
			metadata: metadata,
			Parent:   parent,
		},
	}
}
func skip(reader io.Reader, off int64) (int, error) {
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		return int(n), err
	}
	n, err := io.CopyN(ioutil.Discard, reader, off)
	return int(n), err
}

func (i *File) ReadAt(buff []byte, off int64) (int, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	if len(buff) == 0 {
		return 0, nil
	}
	reader, err := i.Content()
	if err != nil {
		return 0, err
	}
	n, err := skip(reader, off)
	if err != nil {
		return n, err
	}
	num, err := reader.Read(buff)
	return num, err
}

func (i *File) ReadAll() ([]byte, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	reader, err := i.Content()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (i *File) Content() (io.Reader, error) {
	r, err := i.obj.ReadBlob(i.metadata.Hash())
	if err != nil {
		return nil, err
	}
	return r, nil
}

type LimitedWriter struct {
	Buf []byte
	n   int
}

func NewLimitedWriter(n int64) *LimitedWriter {
	return &LimitedWriter{
		Buf: make([]byte, n, n),
	}
}

func (w *LimitedWriter) Write(p []byte) (int, error) {
	l := len(p)
	if w.n+l >= len(w.Buf) {
		copy(w.Buf[w.n:], p)
		return len(w.Buf) - w.n, EDone
	}
	copy(w.Buf[w.n:], p)
	w.n += l
	return l, nil
}

func (w *LimitedWriter) Len() int {
	return len(w.Buf)
}

func (w *LimitedWriter) Read() int {
	return len(w.Buf)
}

var EDone = fmt.Errorf("writer finished")

type DiscardN struct {
	N int
	n int
}

func (w *DiscardN) Write(p []byte) (int, error) {
	l := len(p)
	if w.n+l >= w.N {
		return w.N - w.n, EDone
	}
	w.n += l
	return l, nil
}

type SeqWriter struct {
	writers []io.Writer
	n       int
}

func NewSeqWriter(writers ...io.Writer) *SeqWriter {
	return &SeqWriter{
		writers: writers,
	}
}

type LenWriter struct {
	n int
}

func (w *LenWriter) Write(p []byte) (int, error) {
	w.n += len(p)
	return len(p), nil
}

func (w *LenWriter) Len() int {
	return w.n
}

func (w *SeqWriter) Write(p []byte) (int, error) {
	n, err := w.writers[w.n].Write(p)
	if err != nil && err != EDone {
		return n, err
	}
	if err == nil {
		return n, nil
	}
	// EOF
	for w.n < len(w.writers) {
		w.n++
		nn, err := w.writers[w.n].Write(p[n:])
		n += nn
		if err != nil && err != EDone {
			return n, err
		}
		if err == EDone {
			continue
		}
		if err == nil {
			return n, nil
		}
	}
	return n, err
}

func (i *File) WriteAt(content []byte, offset int64) (n int, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	l := len(content)
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     l,
	}).Debug("SetContent")
	if offset < 0 {
		return 0, e.ENegative
	}
	buf := NewLimitedWriter(offset)
	tail := new(bytes.Buffer)
	w := NewSeqWriter(buf, &DiscardN{N: l}, tail)
	nn, err := i.obj.ReadBlobByWriter(i.metadata.Hash(), w)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if err == io.EOF && int(nn) < buf.Len() {
		return 0, err
	}
	r := io.MultiReader(bytes.NewReader(buf.Buf), bytes.NewReader(content), tail)
	lenWriter := new(LenWriter)
	rr := io.TeeReader(r, lenWriter)
	hash, err := i.obj.WriteBlob(rr)
	if err != nil {
		return 0, err
	}
	i.metadata = i.metadata.Builder().
		Hash(hash).Size(int64(lenWriter.Len())).Build()
	return l, nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	buf := NewLimitedWriter(size)
	if size != 0 {
		_, err := i.obj.ReadBlobByWriter(i.metadata.Hash(), buf)
		if err != nil && err != EDone {
			return err
		}
	}
	hash, err := i.obj.WriteBlob(bytes.NewReader(buf.Buf))
	if err != nil {
		return err
	}
	i.metadata = i.metadata.Builder().
		Hash(hash).Size(size).Build()
	return nil
}

func (i *File) Readdirnames(n int, offset int) (names []string, err error) {
	if i == nil {
		return nil, e.ErrInvalid
	}
	return nil, e.EIsFile
}

func (i *File) Readdir(n int, offset int) ([]*object.Metadata, error) {
	if i == nil {
		return nil, e.ErrInvalid
	}
	return nil, e.EIsFile
}

func (i *File) Close() error {
	err := i.ItemBase.Close()
	if err != nil {
		return err
	}
	return nil
}
