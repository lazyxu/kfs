package object

import (
	"crypto/sha256"
	"io"
	"strings"
	"testing"

	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/storage/memory"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Create a file", t, func() {
		str := "hello, world!"
		hashFunc := func() kfscrypto.Hash {
			return kfscrypto.FromStdHash(sha256.New())
		}
		s := memory.New(hashFunc)
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
