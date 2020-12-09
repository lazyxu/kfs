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
		obj := Init(s)
		Convey("Write to storage", func() {
			key, err1 := obj.WriteBlob(strings.NewReader(str))
			So(err1, ShouldBeNil)
			Convey("Read from storage", func() {
				blob2, err2 := obj.ReadBlob(key)
				So(err2, ShouldBeNil)
				Convey("Should be same", func() {
					buf := new(strings.Builder)
					n, err3 := io.Copy(buf, blob2)
					So(err3, ShouldBeNil)
					So(n, ShouldEqual, len(str))
					So(str, ShouldEqual, buf.String())
				})
			})
		})
	})
}
