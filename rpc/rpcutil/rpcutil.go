package rpcutil

import (
	"encoding/binary"
	"io"
	"strings"

	"google.golang.org/protobuf/proto"
)

func WriteCommandType(w io.Writer, commandType CommandType) error {
	return UnexpectedIfError(binary.Write(w, binary.LittleEndian, commandType))
}

func WriteProto(w io.Writer, m proto.Message) error {
	req, err := proto.Marshal(m)
	if err != nil {
		return UnexpectedIfError(err)
	}
	err = binary.Write(w, binary.LittleEndian, uint64(len(req)))
	if err != nil {
		return UnexpectedIfError(err)
	}
	_, err = w.Write(req)
	if err != nil {
		return UnexpectedIfError(err)
	}
	return nil
}

func ReadProto(r io.Reader, m proto.Message) error {
	var length uint64
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return UnexpectedIfError(err)
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return UnexpectedIfError(err)
	}
	err = proto.Unmarshal(buf, m)
	if err != nil {
		return UnexpectedIfError(err)
	}
	return nil
}

func ReadCommandType(r io.Reader) (commandType CommandType, err error) {
	err = UnexpectedIfError(binary.Read(r, binary.LittleEndian, &commandType))
	return
}

func WriteInvalid(w io.Writer, err error) error {
	err1 := binary.Write(w, binary.LittleEndian, EInvalid)
	if err1 != nil {
		return UnexpectedIfError(err1)
	}
	return WriteString(w, err.Error())
}

func WriteOK(w io.Writer) error {
	return UnexpectedIfError(binary.Write(w, binary.LittleEndian, EOK))
}

func ReadStatus(r io.Reader) (status Status, errMsg string, err error) {
	err = UnexpectedIfError(binary.Read(r, binary.LittleEndian, &status))
	if err != nil {
		return
	}
	if status == EInvalid {
		errMsg, err = ReadString(r)
	}
	return
}

func WriteString(w io.Writer, str string) error {
	length := int64(len(str))
	err := binary.Write(w, binary.LittleEndian, length)
	if err != nil {
		return UnexpectedIfError(err)
	}
	_, err = io.WriteString(w, str)
	return UnexpectedIfError(err)
}

func ReadString(r io.Reader) (string, error) {
	var length int64
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return "", UnexpectedIfError(err)
	}
	var sb strings.Builder
	_, err = io.CopyN(&sb, r, length)
	if err != nil {
		return "", UnexpectedIfError(err)
	}
	return sb.String(), nil
}
