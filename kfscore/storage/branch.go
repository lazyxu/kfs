package storage

import (
	"encoding/json"
	"errors"
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
	Encode(i *Branch, w io.Writer)
	Decode(i *Branch, r io.Reader)
}

var (
	branchEncoderDecoder BranchEncoderDecoder = &branchYamlEncoderDecoder{}
)

type branchYamlEncoderDecoder struct {
}

func (g *branchYamlEncoderDecoder) Encode(i *Branch, w io.Writer) {
	e := yaml.NewEncoder(w)
	err := e.Encode(i)
	defer e.Close()
	if err != nil {
		panic(err)
	}
}

func (g *branchYamlEncoderDecoder) Decode(i *Branch, r io.Reader) {
	e := yaml.NewDecoder(r)
	err := e.Decode(i)
	if err != nil {
		panic(err)
	}
}

type branchJsonEncoderDecoder struct {
}

func (g *branchJsonEncoderDecoder) Encode(i *Branch, w io.Writer) {
	e := json.NewEncoder(w)
	err := e.Encode(i)
	if err != nil {
		panic(err)
	}
}

func (g *branchJsonEncoderDecoder) Decode(i *Branch, r io.Reader) {
	e := json.NewDecoder(r)
	err := e.Decode(i)
	if err != nil {
		panic(err)
	}
}

func (s *Storage) readBranch(branch string, cb func(b *Branch)) {
	p := path.Join(s.root, "branch", branch)
	f, err := os.OpenFile(p, os.O_RDONLY, filePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := &Branch{}
	branchEncoderDecoder.Decode(b, f)
	cb(b)
}

func (s *Storage) GetBranchHash(branch string) string {
	var hash string
	s.readBranch(branch, func(b *Branch) {
		hash = b.BranchHash
	})
	return hash
}

func (s *Storage) CreateBranch(clientID string, branchName string) error {
	p := path.Join(s.root, "branch", branchName)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
			if err != nil {
				panic(err)
			}
			branch := &Branch{
				ClientID:   clientID,
				BranchHash: EmptyDirHash,
			}
			branchEncoderDecoder.Encode(branch, f)
			return nil
		}
		if err != nil {
			panic(err)
		}
	}
	return errors.New("该分支已存在")
}

func (s *Storage) DeleteBranch(clientID string, branchName string) error {
	p := path.Join(s.root, "branch", branchName)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("该分支不存在")
		}
		panic(err)
	}
	f, err := os.OpenFile(p, os.O_RDONLY, filePerm)
	if err != nil {
		panic(err)
	}
	b := &Branch{}
	branchEncoderDecoder.Decode(b, f)
	if b.ClientID != clientID {
		return errors.New("只能删除本客户端创建的分支")
	}
	f.Close()
	return os.Remove(p)
}

func (s *Storage) RenameBranch(clientID string, old string, new string) error {
	oldPath := path.Join(s.root, "branch", old)
	_, err := os.Stat(oldPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("该分支不存在")
		}
		panic(err)
	}
	newPath := path.Join(s.root, "branch", new)
	_, err = os.Stat(newPath)
	if err == nil {
		return errors.New("名称冲突")
	}
	if !os.IsNotExist(err) {
		panic(err)
	}
	f, err := os.OpenFile(oldPath, os.O_RDONLY, filePerm)
	if err != nil {
		panic(err)
	}
	b := &Branch{}
	branchEncoderDecoder.Decode(b, f)
	if b.ClientID != clientID {
		return errors.New("只能重命名本客户端创建的分支")
	}
	f.Close()
	err = os.Rename(oldPath, newPath)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *Storage) ListBranch(cb func(branchName string, ClientID string)) {
	p := path.Join(s.root, "branch")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := path.Join(p, file.Name())
		f, err := os.OpenFile(filePath, os.O_RDONLY, filePerm)
		if err != nil {
			panic(err)
		}
		b := &Branch{}
		branchEncoderDecoder.Decode(b, f)
		cb(file.Name(), b.ClientID)
	}
}
