package local

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestBackup(t *testing.T) {
	os.RemoveAll(testDir)
	ctx := context.Background()
	Visitors := []Visitor{
		&FileSizeVisitor{MaxFileSize: 10},
		&CountVisitor{},
	}
	backupCtx := NewBackupCtx(ctx, "../..", Visitors...)
	done := make(chan struct{})
	tick := time.Tick(time.Second / 2)
	go func() {
		for {
			<-tick
			status := backupCtx.GetStatus()
			fmt.Printf("%+v\n", status)
			if status.Done {
				done <- struct{}{}
			}
		}
	}()
	_, err := backupCtx.Scan()
	if err != nil {
		t.Error(err)
		return
	}
	<-done
	for _, visitor := range Visitors {
		fmt.Printf("%+v\n", visitor)
	}
}
