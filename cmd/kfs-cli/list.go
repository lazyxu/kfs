package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/spf13/cobra"
)

func printHeader(total int) error {
	fmt.Printf("total %d\n", total)
	if total != 0 {
		fmt.Printf("mode      \tcount\ttotalCount\thash\tsize\tmodifyTime         \tname\n")
	}
	return nil
}

func formatCount(mode uint64, count uint64) string {
	if !os.FileMode(mode).IsDir() {
		return strings.Repeat(" ", 5)
	}
	return fmt.Sprintf("%5d", count)
}

func printBody(w io.Writer, dirItem sqlite.IDirItem, isHumanize bool) {
	modifyTime := time.Unix(0, int64(dirItem.GetModifyTime())).Format("2006-01-02 15:04:05")
	if isHumanize {
		fmt.Fprintf(w, "%s\t%s\t     %s\t%s\t%s\t%s\t%s\n",
			os.FileMode(dirItem.GetMode()).String(),
			formatCount(dirItem.GetMode(), dirItem.GetCount()), formatCount(dirItem.GetMode(), dirItem.GetTotalCount()),
			dirItem.GetHash()[:4], humanize.Bytes(dirItem.GetSize()), modifyTime, dirItem.GetName())
	} else {
		fmt.Fprintf(w, "%s\t%s\t     %s\t%s\t%d\t%s\t%s\n",
			os.FileMode(dirItem.GetMode()).String(),
			formatCount(dirItem.GetMode(), dirItem.GetCount()), formatCount(dirItem.GetMode(), dirItem.GetTotalCount()),
			dirItem.GetHash()[:4], dirItem.GetSize(), modifyTime, dirItem.GetName())
	}
}

func listCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "ls",
		Example: "kfs-cli ls .",
		Args:    cobra.RangeArgs(0, 1),
		Run:     runList,
	}
	cmd.PersistentFlags().Bool(HumanizeStr, true, "")
	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

	p := ""
	if len(args) != 0 {
		p = args[0]
	}
	isHumanize := cmd.Flag(HumanizeStr).Value.String() == "true"

	err = fs.List(cmd.Context(), branchName, p, func(total int64) error {
		cmd.Printf("total %d\n", total)
		if total != 0 {
			cmd.Printf("mode      \tcount\ttotalCount\thash\tsize\tmodifyTime         \tname\n")
		}
		return nil
	}, func(item sqlite.IDirItem) error {
		printBody(cmd.OutOrStdout(), item, isHumanize)
		return nil
	})
}
