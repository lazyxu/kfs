package storage

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type Directory = []*Metadata

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
	Encode(i Directory) []byte
	Decode(data []byte) Directory
}

var (
	DefaultDirectoryEncoderDecoder DirectoryEncoderDecoder = &directoryYamlEncoderDecoder{}
)

type directoryYamlEncoderDecoder struct {
}

func (g *directoryYamlEncoderDecoder) Encode(i Directory) []byte {
	data, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

func (g *directoryYamlEncoderDecoder) Decode(data []byte) Directory {
	var i Directory
	err := yaml.Unmarshal(data, &i)
	if err != nil {
		panic(err)
	}
	return i
}

type directoryJsonEncoderDecoder struct {
}

func (g *directoryJsonEncoderDecoder) Encode(i Directory) []byte {
	data, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

func (g *directoryJsonEncoderDecoder) Decode(data []byte) Directory {
	var i Directory
	err := json.Unmarshal(data, &i)
	if err != nil {
		panic(err)
	}
	return i
}
