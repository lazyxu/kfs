package dirignore

import (
	"path"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/sirupsen/logrus"
)

type DirGitIgnore struct {
	Path   string
	Ignore *ignore.GitIgnore
	Parent *DirGitIgnore
}

func New(dir string) *DirGitIgnore {
	gitIgnore, err := ignore.CompileIgnoreLines(".git")
	if err != nil {
		logrus.WithError(err).Error("compileIgnore")
	}
	root := &DirGitIgnore{
		Path:   dir,
		Ignore: gitIgnore,
	}
	return root.Enter(dir)
}

func (i *DirGitIgnore) Enter(dir string) *DirGitIgnore {
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
	return &DirGitIgnore{
		Path:   dir,
		Ignore: gitIgnore,
		Parent: i,
	}
}

func (i *DirGitIgnore) Exit(dir string) *DirGitIgnore {
	//return i.Parent
	if i.Path == dir {
		return i.Parent
	}
	return i
}
