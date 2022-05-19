package initialization

import (
	"path/filepath"
)

func local(addr string, branchName string, description string) (string, error) {
	var err error
	addr, err = filepath.Abs(addr)
	if err != nil {
		return "", err
	}
	// TODO: create branch with description
	return addr, nil
}
