package storage

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Directory struct {
	Items []*Metadata `yaml:"items"`
}

type Metadata struct {
	Mode       os.FileMode `yaml:"mode"`
	BirthTime  int64       `yaml:"birth_time"`
	ModifyTime int64       `yaml:"modify_time"`
	ChangeTime int64       `yaml:"change_time"`
	Size       int64       `yaml:"size"`
	Name       string      `yaml:"name"`
	Hash       string      `yaml:"hash"`
}

type DirectoryEncoderDecoder interface {
	Encode(i *Directory, w io.Writer)
	Decode(i *Directory, r io.Reader)
}

var (
	directoryEncoderDecoder DirectoryEncoderDecoder = &directoryYamlEncoderDecoder{}
)

type directoryYamlEncoderDecoder struct {
}

func (g *directoryYamlEncoderDecoder) Encode(i *Directory, w io.Writer) {
	e := yaml.NewEncoder(w)
	err := e.Encode(i)
	defer e.Close()
	if err != nil {
		panic(err)
	}
}

func (g *directoryYamlEncoderDecoder) Decode(i *Directory, r io.Reader) {
	e := yaml.NewDecoder(r)
	err := e.Decode(i)
	if err != nil {
		panic(err)
	}
}
