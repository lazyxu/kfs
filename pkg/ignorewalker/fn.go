package ignorewalker

import (
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/lazyxu/kfs/pkg/ignorewalker/dirignore"
	"github.com/sirupsen/logrus"
)

func walkFn(p string, ignore *dirignore.DirGitIgnore, info os.FileInfo, err error) error {
	if info == nil || err != nil || ignore == nil {
		return nil
	}
	size := uint64(info.Size())
	total++
	totalSize += size
	i := ignore
	for {
		for {
			if i == nil {
				goto exit
			}
			if i.Ignore != nil {
				break
			}
			i = i.Parent
		}
		rel, err := filepath.Rel(i.Path, p)
		if err != nil {
			logrus.WithError(err).Error("rel")
			return err
		}
		match := i.Ignore.MatchesPath(rel)
		logrus.WithFields(logrus.Fields{
			"path":   i.Path,
			"ignore": match,
		}).Trace("ignore")
		if match {
			logrus.WithFields(logrus.Fields{
				"total":          total,
				"totalSize":      humanize.Bytes(totalSize),
				"notIgnored":     notIgnored,
				"notIgnoredSize": humanize.Bytes(notIgnoredSize),
			}).Trace(p)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		i = i.Parent
	}
exit:
	if info.IsDir() {
		return nil
	}
	ignore.Files = append(ignore.Files, &dirignore.File{
		Path: p,
		Size: size,
	})
	ignore.Size += size
	notIgnored++
	notIgnoredSize += size
	logrus.WithFields(logrus.Fields{
		"total":          total,
		"totalSize":      humanize.Bytes(totalSize),
		"notIgnored":     notIgnored,
		"notIgnoredSize": humanize.Bytes(notIgnoredSize),
	}).Info(p)
	return nil
}
