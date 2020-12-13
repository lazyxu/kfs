package object

import (
	"crypto/sha256"
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
				buf := new(strings.Builder)
				_, err2 := obj.ReadBlobByWriter(key, buf)
				So(err2, ShouldBeNil)
				Convey("Should be same", func() {
					So(str, ShouldEqual, buf.String())
				})
			})
		})
	})
}
