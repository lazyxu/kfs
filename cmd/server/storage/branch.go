package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Branch struct {
	ClientID   string `yaml:"client_id"`
	BranchHash string `yaml:"branch_hash"`
}

type BranchEncoderDecoder interface {
	Encode(i *Branch, w io.Writer) (err error)
	Decode(i *Branch, r io.Reader) (err error)
}

var (
	branchEncoderDecoder BranchEncoderDecoder = &branchYamlEncoderDecoder{}
)

type branchYamlEncoderDecoder struct {
}

func (g *branchYamlEncoderDecoder) Encode(i *Branch, w io.Writer) error {
	e := yaml.NewEncoder(w)
	err := e.Encode(i)
	defer e.Close()
	if err != nil {
		return err
	}
	return nil
}

func (g *branchYamlEncoderDecoder) Decode(d *Branch, r io.Reader) error {
	e := yaml.NewDecoder(r)
	err := e.Decode(d)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateBranch(clientID string, branchName string) error {
	p := path.Join(s.root, "branch", branchName)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
			if err != nil {
				return err
			}
			branch := &Branch{
				ClientID:   clientID,
				BranchHash: EmptyDirHash,
			}
			err = branchEncoderDecoder.Encode(branch, f)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (s *Storage) ListBranch(cb func(branchName string, ClientID string) error) error {
	p := path.Join(s.root, "branch")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := path.Join(p, file.Name())
		f, err := os.OpenFile(filePath, os.O_RDONLY, filePerm)
		if err != nil {
			return err
		}
		b := &Branch{}
		err = branchEncoderDecoder.Decode(b, f)
		if err != nil {
			return err
		}
		err = cb(file.Name(), b.ClientID)
		if err != nil {
			return err
		}
	}
	return nil
}
