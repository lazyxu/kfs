package rootdirectory

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"

	"github.com/lazyxu/kfs/kfscore/object"

	"github.com/lazyxu/kfs/warpper/grpcweb/pb"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/sirupsen/logrus"
)

const (
	MethodInvalid int32 = iota
	MethodUploadBlob
	MethodUploadTree
	MethodUpdateBranch
)

const (
	CompressNone int32 = iota
	CompressZlib
)

func process(conn net.Conn, s storage.Storage) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	var headerLen uint64
	err := binary.Read(reader, binary.LittleEndian, &headerLen)
	if err == io.EOF {
		return
	}
	if err != nil {
		logrus.WithError(err).Error("method")
		return
	}
	logrus.Infoln("headerLen", headerLen)

	rawHeader := make([]byte, headerLen)
	n, err := reader.Read(rawHeader)
	if err != nil {
		logrus.WithError(err).Error("branch")
		return
	}
	if uint64(n) != headerLen {
		logrus.WithField("n", n).Error("header len not enough")
		return
	}

	var header pb.Header
	err = proto.Unmarshal(rawHeader, &header)
	if err != nil {
		logrus.WithError(err).Error("branch")
		return
	}
	logrus.Infoln("header", header.String())

	if header.Method == MethodUploadBlob {
		obj := object.Init(s)
		if len(header.Hash) != 0 {
			exist, err := obj.S.Exist(storage.TypBlob, header.Hash)
			if err != nil {
				logrus.WithError(err).Error("Exist")
				return
			}
			if exist {
				_, err = conn.Write([]byte{1})
			} else {
				_, err = conn.Write([]byte{0})
			}
			if err != nil {
				logrus.WithError(err).Error("Write")
				return
			}
		}
		r := io.LimitReader(reader, int64(header.RawSize))
		hash, err := obj.WriteBlob(r)
		if err != nil {
			logrus.WithError(err).Error("WriteBlob")
			return
		}
		logrus.Infoln("hash", hash)

		hashLen := uint8(len(hash))
		err = binary.Write(conn, binary.LittleEndian, hashLen)
		if err != nil {
			logrus.WithError(err).Error("Write hashLen")
			return
		}
		_, err = conn.Write([]byte(hash))
		if err != nil {
			logrus.WithError(err).Error("Write")
			return
		}
	} else if header.Method == MethodUploadTree {
		obj := object.Init(s)
		r := io.LimitReader(reader, int64(header.RawSize))
		buf := make([]byte, header.RawSize)
		_, err = r.Read(buf)
		if err != nil {
			logrus.WithError(err).Error("Read")
			return
		}
		infos := &pb.FileInfos{}
		err = proto.Unmarshal(buf, infos)
		if err != nil {
			logrus.WithError(err).Error("Unmarshal")
			return
		}
		t := obj.NewTree()
		for _, info := range infos.Info {
			if info == nil || info.Type == "" {
				continue
			}
			var item *object.Metadata
			if info.Type == "file" {
				item = obj.NewFileMetadata(info.Name, os.FileMode(info.Mode)).Builder().
					Hash(info.Hash).Size(info.Size).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
			} else if info.Type == "dir" {
				item = obj.NewDirMetadata(info.Name, os.FileMode(info.Mode)).Builder().
					Hash(info.Hash).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
			}
			t.Items = append(t.Items, item)
		}
		hash, err := obj.WriteTree(t)
		if err != nil {
			logrus.WithError(err).Error("WriteTree")
			return
		}
		logrus.Infoln("hash", hash)

		hashLen := uint8(len(hash))
		err = binary.Write(conn, binary.LittleEndian, hashLen)
		if err != nil {
			logrus.WithError(err).Error("Write hashLen")
			return
		}
		_, err = conn.Write([]byte(hash))
		if err != nil {
			logrus.WithError(err).Error("Write")
			return
		}
	}
}

func Socket(s storage.Storage, port int) {
	listen, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}
	logrus.WithField("port", port).Info("Listening socket")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		go process(conn, s)
	}
}
