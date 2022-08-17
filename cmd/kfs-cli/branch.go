package main

import (
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func branchCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "branch",
	}

	cmd.AddCommand(branchCheckoutCmd())
	cmd.AddCommand(branchInfoCmd())
	cmd.AddCommand(branchUpdateCmd())
	cmd.PersistentFlags().String(DescriptionStr, "", "branch description")
	return cmd
}

func branchCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "checkout",
		Example: "kfs-cli branch checkout branchName",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runCheckoutBranch,
	}
}

func checkoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "checkout",
		Example: "kfs-cli checkout branchName",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runCheckoutBranch,
	}
}

func runCheckoutBranch(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, _, _ := loadFs(cmd)

	branchName := args[0]

	_, err = fs.Checkout(cmd.Context(), branchName)
	if err != nil {
		return
	}
	cmd.Printf("switch to branch '%s'\n", branchName)
	viper.Set(BranchNameStr, branchName)
	err = viper.WriteConfig()
}

func branchInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Example: "kfs-cli branch info branchName",
		Args:    cobra.RangeArgs(0, 1),
		Run:     runBranchInfo,
	}
}

func runBranchInfo(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

	if len(args) != 0 {
		branchName = args[0]
	}

	branch, err := fs.BranchInfo(cmd.Context(), branchName)
	if err != nil {
		return
	}

	cmd.Printf("description: %s\n", branch.GetDescription())
	cmd.Printf("commitId: %d\n", branch.GetCommitId())
	cmd.Printf("size: %d\n", branch.GetSize())
	cmd.Printf("count: %d\n", branch.GetCount())
}

func branchUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		Example: "kfs-cli branch update branchName",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runCheckoutBranch,
	}
}
