package object

import (
	"io"
	"strings"
	"testing"

	"github.com/lazyxu/kfs/storage/memory"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Create a file", t, func() {
		str := "hello, world!"
		s := memory.New()
		blob1 := &Blob{Reader: strings.NewReader(str)}
		Convey("Write to storage", func() {
			key, err1 := blob1.Write(s)
			So(err1, ShouldBeNil)
			Convey("Read from storage", func() {
				blob2 := new(Blob)
				err2 := blob2.Read(s, key)
				So(err2, ShouldBeNil)
				Convey("Should be same", func() {
					buf := new(strings.Builder)
					n, err3 := io.Copy(buf, blob2.Reader)
					So(err3, ShouldBeNil)
					So(n, ShouldEqual, len(str))
					So(str, ShouldEqual, buf.String())
				})
			})
		})
	})
}
