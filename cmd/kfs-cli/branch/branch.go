package branch

import (
	"context"
	"fmt"
	"os"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"

	"github.com/lazyxu/kfs/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Cmd = &cobra.Command{
	Use: "branch",
}

func init() {
	branchCheckoutCmd.PersistentFlags().String(DescriptionStr, "", "branch description")
	Cmd.AddCommand(branchCheckoutCmd)
	Cmd.AddCommand(branchInfoCmd)
}

var branchCheckoutCmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli branch checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     RunCheckoutBranch,
}

func RunCheckoutBranch(cmd *cobra.Command, args []string) {
	var err error
	defer ExitWithError(err)
	remoteAddr := viper.GetString(ServerAddrStr)
	oldBranchName := viper.GetString(BranchNameStr)
	description := cmd.Flag(DescriptionStr).Value.String()
	branchName := args[0]
	fmt.Printf("remoteAddr=%s\n", remoteAddr)
	fmt.Printf("branch=%s\n", oldBranchName)
	exist, err := Checkout(remoteAddr, branchName, description)
	if exist {
		return
	}
	viper.Set(BranchNameStr, branchName)
	err = viper.WriteConfig()
}

func Checkout(remoteAddr string, branchName string, description string) (bool, error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return true, nil
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	ctx := context.Background()
	_, err = c.BranchCheckout(ctx, &pb.BranchReq{
		BranchName:  branchName,
		Description: description,
	})
	if err != nil {
		return true, nil
	}
	fmt.Printf("switch to branch '%s'\n", branchName)
	return false, err
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	remoteAddr := viper.GetString(ServerAddrStr)
	var branchName string
	if len(args) != 0 {
		branchName = args[0]
	} else {
		branchName = viper.GetString(BranchNameStr)
	}
	fmt.Printf("remoteAddr=%s\n", remoteAddr)
	fmt.Printf("branch=%s\n", branchName)
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	ctx := context.Background()
	branch, err := c.BranchInfo(ctx, &pb.BranchInfoReq{
		BranchName: branchName,
	})
	if err != nil {
		return
	}
	fmt.Printf("description: %s\n", branch.Description)
	fmt.Printf("commitId: %d\n", branch.CommitId)
	fmt.Printf("size: %d\n", branch.Size)
	fmt.Printf("count: %d\n", branch.Count)
}
