package dirignore

import (
	"path"
	"runtime"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/sirupsen/logrus"
)

type File struct {
	Path string
	Size uint64
}

type DirGitIgnore struct {
	Path     string
	Ignore   *ignore.GitIgnore
	Parent   *DirGitIgnore
	Children []*DirGitIgnore
	Root     *DirGitIgnore
	Files    []*File
	DirSize  map[string]uint64
	Size     uint64
}

var defaultIgnore = []string{
	".git",
	".npm",
	".node-gyp",
	".vscode",
	".nvm",
	"go/pkg",
	"go/bin",
	".Trash",
	".gradle",
	".dropbox",
}

var defaultIgnoreMac = append(defaultIgnore, "Library",
	"Applications",
	".DS_Store",
	"*.app",
)

func New(dir string) *DirGitIgnore {
	var ignores []string
	switch runtime.GOOS {
	case "darwin":
		ignores = defaultIgnoreMac
	default:
		ignores = defaultIgnore
		return nil
	}
	gitIgnore, err := ignore.CompileIgnoreLines(ignores...)
	if err != nil {
		logrus.WithError(err).Error("compileIgnore")
	}
	osRoot := &DirGitIgnore{
		Path:   dir,
		Ignore: gitIgnore,
	}
	root := osRoot.Enter(dir)
	root.Root = root
	return root
}

func (i *DirGitIgnore) Enter(dir string) *DirGitIgnore {
	// TODO: do not ignore files?
	gitIgnore, err := ignore.CompileIgnoreFile(path.Join(dir, ".gitignore"))
	if err != nil {
		//return &DirGitIgnore{
		//	Path:   dir,
		//	Ignore: gitIgnore,
		//	Parent: i,
		//}
		return i
	}
	logrus.WithFields(logrus.Fields{
		"path":      path.Join(dir, ".gitignore"),
		"gitIgnore": gitIgnore,
	}).Debug("enter")
	newDir := &DirGitIgnore{
		Path:     dir,
		Ignore:   gitIgnore,
		Children: []*DirGitIgnore{},
		Files:    []*File{},
		Parent:   i,
		Root:     i.Root,
	}
	i.Children = append(i.Children, newDir)
	return newDir
}

func (i *DirGitIgnore) Exit(dir string) *DirGitIgnore {
	//return i.Parent
	if i.Path == dir {
		return i.Parent
	}
	return i
}

func (i *DirGitIgnore) CalcDirSize() {
	// TODO: recursive?
	i.DirSize = make(map[string]uint64)
	if len(i.Files) == 0 {
		return
	}
	for _, file := range i.Files {
		dirName := path.Dir(file.Path)
		if _, ok := i.DirSize[dirName]; ok {
			i.DirSize[dirName] += file.Size
		} else {
			i.DirSize[dirName] = file.Size
		}
	}
}
