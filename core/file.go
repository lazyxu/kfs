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
	offset int64 // file pointer offset
	closed bool  // set if handle has been closed
	opened bool
}

func NewFile(kfs *KFS, name string) *File {
	return &File{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: object.NewFileMetadata(name),
		},
	}
}

func (i *File) Read(buff []byte) (int, error) {
	n, err := i.ReadAt(buff, i.offset)
	if err != nil {
		return 0, err
	}
	i.offset += int64(n)
	return n, nil
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
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		if err != nil {
			return int(n), err
		}
	default:
		n, err := io.CopyN(ioutil.Discard, r, off)
		if err != nil {
			return int(n), err
		}
	}
	num, err := reader.Read(buff)
	return num, err
}

func (i *File) getContent() (io.Reader, error) {
	blob := new(object.Blob)
	err := blob.Read(i.kfs.scheduler, i.Metadata.Hash)
	if err != nil {
		return nil, err
	}
	return blob.Reader, nil
}

func (i *File) Write(content []byte) (n int, err error) {
	n, err = i.WriteAt(content, i.offset)
	if err != nil {
		return 0, err
	}
	i.offset += int64(n)
	return n, nil
}

func (i *File) WriteAt(content []byte, offset int64) (n int, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     len(content),
	}).Debug("SetContent")
	buf := make([]byte, offset)
	blob := new(object.Blob)
	err = blob.Read(i.kfs.scheduler, i.Metadata.Hash)
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
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.kfs.scheduler)
	if err != nil {
		return 0, err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = int64(len(content))
	return int(i.Metadata.Size), nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	blob := new(object.Blob)
	if size != 0 {
		err := blob.Read(i.kfs.scheduler, i.Metadata.Hash)
		if err != nil {
			return err
		}
		_, err = blob.Reader.Read(content)
		if err != nil {
			return err
		}
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.kfs.scheduler)
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
	i.offset = 0
	i.closed = true
	i.opened = false
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
func (i *File) Open(flags int) (fd Handle, err error) {
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

	if read && write {
		fd, err = i.openRW(flags)
	} else if write {
		fd, err = i.openWrite(flags)
	} else if read {
		fd, err = i.openRead()
	}
	// TODO: create
	return fd, err
}

// openRead open the file for read
func (f *File) openRead() (fh *ReadFileHandle, err error) {
	return newReadFileHandle(f.kfs, f.Path()), nil
}

// openWrite open the file for write
func (f *File) openWrite(flags int) (fh *WriteFileHandle, err error) {
	return newWriteFileHandle(f.kfs, f.Path()), nil
}

// openRW open the file for read and write using a temporay file
//
// It uses the open flags passed in.
func (f *File) openRW(flags int) (fh *RWFileHandle, err error) {
	return newRWFileHandle(f.kfs, f.Path()), nil
}
