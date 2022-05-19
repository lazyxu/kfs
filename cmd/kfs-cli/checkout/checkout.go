package checkout

import (
	"github.com/lazyxu/kfs/cmd/kfs-cli/branch"
	"github.com/lazyxu/kfs/cmd/kfs-cli/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     branch.RunCheckoutBranch,
}

func init() {
	Cmd.PersistentFlags().String(utils.DescriptionStr, "", "branch description")
}
