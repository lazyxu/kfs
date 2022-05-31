package local

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestBackup(t *testing.T) {
	os.RemoveAll(testDir)
	ctx := context.Background()
	Visitors := []Visitor[any]{
		&FileSizeVisitor{MaxFileSize: 10},
		&CountVisitor{},
	}
	walker := NewWalker(ctx, "../..", Visitors...)
	_, err := walker.Walk(false)
	if err != nil {
		t.Error(err)
		return
	}
	for _, visitor := range Visitors {
		fmt.Printf("%+v\n", visitor)
	}
}
