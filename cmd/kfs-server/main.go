package main

import (
	"fmt"
	"github.com/lazyxu/kfs/rpc/server"
	"os"
)

func main() {
	{
		//f, err := os.OpenFile("G:\\备份\\我的设备\\iPhone 13 Pro\\Internal Storage\\DCIM原图\\202304__\\IMG_7547.HEIC", os.O_RDONLY, 0o200)
		f, err := os.OpenFile("G:\\备份\\我的设备\\iPhone 13 Pro\\Internal Storage\\DCIM原图\\202304__\\IMG_7548.MOV", os.O_RDONLY, 0o200)
		if err != nil {
			return
		}
		defer f.Close()
		e, err := server.GetExifDataWithReadAtSeeker(f)
		if err != nil {
			fmt.Printf("GetExifDataWithReadAtSeeker: %+v\n", e)
			return
		}
		fmt.Printf("GetExifDataWithReadAtSeeker: %+v\n", e)
		return
	}
	AnalysisExifProcess()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
