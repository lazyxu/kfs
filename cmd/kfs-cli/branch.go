package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use: "branch",
}

func init() {
	branchCmd.AddCommand(branchCheckoutCmd)
	branchCmd.AddCommand(branchInfoCmd)
	branchCmd.AddCommand(branchUpdateCmd)
	branchUpdateCmd.PersistentFlags().String(DescriptionStr, "", "branch description")
}

var branchCheckoutCmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli branch checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}

var checkoutCmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}

func runCheckoutBranch(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(err)
	}()
	serverType := viper.GetString(ServerTypeStr)
	serverAddr := viper.GetString(ServerAddrStr)
	oldBranchName := viper.GetString(BranchNameStr)
	fmt.Printf("%s: %s\n", ServerTypeStr, serverType)
	fmt.Printf("%s: %s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s: %s\n", BranchNameStr, oldBranchName)

	branchName := args[0]

	fs, err := getFS(serverType, serverAddr)
	if err != nil {
		return
	}
	defer fs.Close()

	_, err = fs.Checkout(cmd.Context(), branchName)
	if err != nil {
		return
	}
	fmt.Printf("switch to branch '%s'\n", branchName)
	viper.Set(BranchNameStr, branchName)
	err = viper.WriteConfig()
}

var branchInfoCmd = &cobra.Command{
	Use:     "info",
	Example: "kfs-cli branch info branchName",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runBranchInfo,
}

func runBranchInfo(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(err)
	}()
	serverType := viper.GetString(ServerTypeStr)
	serverAddr := viper.GetString(ServerAddrStr)
	var branchName string
	if len(args) != 0 {
		branchName = args[0]
	} else {
		branchName = viper.GetString(BranchNameStr)
	}
	fmt.Printf("%s: %s\n", ServerTypeStr, serverType)
	fmt.Printf("%s: %s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s: %s\n", BranchNameStr, branchName)

	fs, err := getFS(serverType, serverAddr)
	if err != nil {
		return
	}
	defer fs.Close()

	branch, err := fs.BranchInfo(cmd.Context(), branchName)
	if err != nil {
		return
	}

	fmt.Printf("description: %s\n", branch.GetDescription())
	fmt.Printf("commitId: %d\n", branch.GetCommitId())
	fmt.Printf("size: %d\n", branch.GetSize())
	fmt.Printf("count: %d\n", branch.GetCount())
}

var branchUpdateCmd = &cobra.Command{
	Use:     "update",
	Example: "kfs-cli branch update branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}
