package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/object"

	"github.com/sirupsen/logrus"
)

type File struct {
	ItemBase
}

func NewFile(kfs *KFS, name string) *File {
	return &File{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: kfs.baseObject.NewFileMetadata(name),
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
	reader, err := i.getContent()
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
	reader, err := i.getContent()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (i *File) getContent() (io.Reader, error) {
	blob := new(object.Blob)
	err := blob.Read(i.kfs.storage, i.Metadata.Hash)
	if err != nil {
		return nil, err
	}
	return blob.Reader, nil
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
	buf := make([]byte, offset)
	blob := new(object.Blob)
	err = blob.Read(i.kfs.storage, i.Metadata.Hash)
	if err != nil {
		return 0, err
	}
	if offset != 0 {
		_, err = blob.Reader.Read(buf)
		if err != nil {
			return 0, err
		}
	}
	content = append(buf, content...)
	n, err = skip(blob.Reader, int64(l))
	if err != nil && err != io.EOF {
		return n, err
	}
	if err != io.EOF {
		remain, err := ioutil.ReadAll(blob.Reader)
		if err != nil {
			return 0, err
		}
		content = append(content, remain...)
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.kfs.storage)
	if err != nil {
		return 0, err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = int64(len(content))
	return l, nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	blob := new(object.Blob)
	if size != 0 {
		err := blob.Read(i.kfs.storage, i.Metadata.Hash)
		if err != nil {
			return err
		}
		_, err = blob.Reader.Read(content)
		if err != nil {
			return err
		}
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.kfs.storage)
	if err != nil {
		return err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = size
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

// Open a file according to the flags provided
//
//   O_RDONLY open the file read-only.
//   O_WRONLY open the file write-only.
//   O_RDWR   open the file read-write.
//
//   O_APPEND append data to the file when writing.
//   O_CREATE create a new file if none exists.
//   O_EXCL   used with O_CREATE, file must not exist
//   O_SYNC   open for synchronous I/O.
//   O_TRUNC  if possible, truncate file when opene
//
// We ignore O_SYNC and O_EXCL
func (i *File) Open(flags int) (fd *Handle, err error) {
	var (
		write    bool // if set need write support
		read     bool // if set need read support
		rdwrMode = flags & accessModeMask
	)

	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if rdwrMode == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		return nil, e.ErrInvalid
	}

	if flags&os.O_TRUNC != 0 {
		err := i.Truncate(0)
		if err != nil {
			return nil, err
		}
	}
	// Figure out the read/write intents
	switch {
	case rdwrMode == os.O_RDONLY:
		read = true
	case rdwrMode == os.O_WRONLY:
		write = true
	case rdwrMode == os.O_RDWR:
		read = true
		write = true
	default:
		logrus.Debug(i.Name(), "Can't figure out how to open with flags: 0x%X", flags)
		return nil, e.ErrPermission
	}

	return &Handle{
		kfs:    i.kfs,
		path:   i.Path(),
		read:   read,
		write:  write,
		append: flags&os.O_APPEND != 0,
	}, nil
}
