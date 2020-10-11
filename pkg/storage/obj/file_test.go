package obj

import (
	"io"
	"strings"
	"testing"

	"github.com/lazyxu/kfs/storage/memory"
	"github.com/lazyxu/kfs/storage/scheduler"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Create a file", t, func() {
		str := "hello, world!"
		s := scheduler.New(memory.New())
		file1 := &File{Reader: strings.NewReader(str)}
		Convey("Write to storage", func() {
			key, err1 := file1.Write(s)
			So(err1, ShouldBeNil)
			Convey("Read from storage", func() {
				file2 := new(File)
				err2 := file2.Read(s, key)
				So(err2, ShouldBeNil)
				Convey("Should be same", func() {
					buf := new(strings.Builder)
					n, err3 := io.Copy(buf, file2.Reader)
					So(err3, ShouldBeNil)
					So(n, ShouldEqual, len(str))
					So(str, ShouldEqual, buf.String())
				})
			})
		})
	})
}
