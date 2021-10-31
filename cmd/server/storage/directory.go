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
	Encode(i *Directory, w io.Writer) (err error)
	Decode(i *Directory, r io.Reader) (err error)
}

var (
	directoryEncoderDecoder DirectoryEncoderDecoder = &directoryYamlEncoderDecoder{}
)

type directoryYamlEncoderDecoder struct {
}

func (g *directoryYamlEncoderDecoder) Encode(i *Directory, w io.Writer) error {
	e := yaml.NewEncoder(w)
	err := e.Encode(i)
	defer e.Close()
	if err != nil {
		return err
	}
	return nil
}

func (g *directoryYamlEncoderDecoder) Decode(d *Directory, r io.Reader) error {
	e := yaml.NewDecoder(r)
	err := e.Decode(d)
	if err != nil {
		return err
	}
	return nil
}
