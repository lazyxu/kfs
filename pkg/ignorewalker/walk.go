package ignorewalker

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/lazyxu/kfs/pkg/ignorewalker/dirignore"
)

var (
	total          = 0
	totalSize      uint64
	notIgnored     = 0
	notIgnoredSize uint64
)

type WalkFunc func(path string, ignore *dirignore.DirGitIgnore, info os.FileInfo, err error) error

func Walk(root string) (*dirignore.DirGitIgnore, error) {
	dir, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	ignore := dirignore.New(dir)
	return ignore, iWalk(dir, ignore, walkFn)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very
// large directories Walk can be inefficient.
// Walk does not follow symbolic links.
func iWalk(root string, ignore *dirignore.DirGitIgnore, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = walkFn(root, ignore, nil, err)
	} else {
		err = walk(root, ignore, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

// walk recursively descends path, calling walkFn.
func walk(path string, ignore *dirignore.DirGitIgnore, info os.FileInfo, walkFn WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, ignore, info, nil)
	}

	names, err := readDirNames(path)
	err1 := walkFn(path, ignore, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, ignore, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			ignore = ignore.Enter(filename)
			err = walk(filename, ignore, fileInfo, walkFn)
			ignore = ignore.Exit(filename)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
