package initialization

import (
	"github.com/lazyxu/kfs/cmd/kfs-cli/branch"
)

func remote(addr string, branchName string, description string) error {
	_, err := branch.Checkout(addr, branchName, description)
	if err != nil {
		return err
	}
	return nil
}
