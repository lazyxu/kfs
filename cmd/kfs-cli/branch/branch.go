package branch

import (
	"context"
	"fmt"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lazyxu/kfs/core"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "branch",
}

func init() {
	Cmd.AddCommand(branchCheckoutCmd)
	Cmd.AddCommand(branchInfoCmd)
	Cmd.AddCommand(branchUpdateCmd)
	branchUpdateCmd.PersistentFlags().String(DescriptionStr, "", "branch description")
}

var branchCheckoutCmd = &cobra.Command{
	Use:     "checkout",
	Example: "kfs-cli branch checkout branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}

var CheckoutCmd = &cobra.Command{
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

	switch serverType {
	case ServerTypeLocal:
		_, err = core.Checkout(cmd.Context(), serverAddr, branchName)
	case ServerTypeRemote:
		_, err = RemoteCheckout(cmd.Context(), serverAddr, branchName)
	default:
		err = InvalidServerType
	}

	if err != nil {
		return
	}
	fmt.Printf("switch to branch '%s'\n", branchName)
	viper.Set(BranchNameStr, branchName)
	err = viper.WriteConfig()
}

func RemoteCheckout(ctx context.Context, remoteAddr string, branchName string) (bool, error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	resp, err := c.BranchCheckout(ctx, &pb.BranchReq{
		BranchName: branchName,
	})
	return resp.Exist, err
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

	var branch sqlite.IBranch
	switch serverType {
	case ServerTypeLocal:
		branch, err = core.BranchInfo(cmd.Context(), serverAddr, branchName)
	case ServerTypeRemote:
		branch, err = remoteBranchInfo(cmd.Context(), serverAddr, branchName)
	default:
		err = InvalidServerType
	}

	if err != nil {
		return
	}
	fmt.Printf("description: %s\n", branch.GetDescription())
	fmt.Printf("commitId: %d\n", branch.GetCommitId())
	fmt.Printf("size: %d\n", branch.GetSize())
	fmt.Printf("count: %d\n", branch.GetCount())
}

func remoteBranchInfo(ctx context.Context, addr string, branchName string) (sqlite.IBranch, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	return c.BranchInfo(ctx, &pb.BranchInfoReq{
		BranchName: branchName,
	})
}

var branchUpdateCmd = &cobra.Command{
	Use:     "update",
	Example: "kfs-cli branch update branchName",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCheckoutBranch,
}
