package localfs

import (
	"os"
)

type Info interface {
}

type RegularFile struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type Dir struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func TransInfo(dirname string, infos []os.FileInfo) map[string]interface{} {
	var newInfos []Info
	for _, info := range infos {
		var newInfo Info
		if info.IsDir() {
			newInfo = Dir{
				Name: info.Name(),
				Type: "dir",
			}
		} else if info.Mode()&os.ModeType == 0 {
			newInfo = RegularFile{
				Name: info.Name(),
				Type: "file",
				Size: info.Size(),
			}
		} else {
			continue
		}
		newInfos = append(newInfos, newInfo)
	}
	if newInfos == nil {
		newInfos = []Info{}
	}
	return map[string]interface{}{
		"Dirname": dirname,
		"files":   newInfos,
	}
}

type Result struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result"`
}
