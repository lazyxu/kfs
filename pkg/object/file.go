package object

import "crypto/sha256"

type File struct {
	baseObject
	content string
}

func NewEmptyFile(name string) *File {
	return &File{
		baseObject: baseObject{
			TimeImpl: NewTimeImpl(),
			name:     name,
			hash:     EmptyFileHash,
			mode:     DefaultFileMode,
		},
		content: "",
	}
}

func (o *File) Content() string {
	return o.content
}

func (o *File) SetContent(content string) {
	o.content = content
	o.size = int64(len(content))
}

func (o *File) Hash() string {
	hash := sha256.New()
	hash.Write([]byte("file"))
	return string(hash.Sum([]byte(o.content)))
}

func (o File) Clone() Object {
	return &o
}
