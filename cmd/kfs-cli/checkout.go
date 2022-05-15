package main

import "github.com/spf13/cobra"

var checkoutCmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}

func init() {
	checkoutCmd.PersistentFlags().String(descriptionStr, "", "branch description")
}
