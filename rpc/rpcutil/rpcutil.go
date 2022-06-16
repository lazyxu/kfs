package rpcutil

import (
	"encoding/binary"
	"io"
	"strings"
)

func WriteCommandType(w io.Writer, commandType CommandType) error {
	return binary.Write(w, binary.LittleEndian, commandType)
}

func ReadCommandType(r io.Reader) (commandType CommandType, err error) {
	err = binary.Read(r, binary.LittleEndian, &commandType)
	return
}

func WriteErrorExit(w io.Writer, err error) error {
	err1 := binary.Write(w, binary.LittleEndian, EInvalid)
	if err1 != nil {
		return err1
	}
	return WriteString(w, err.Error())
}

func WriteSuccessExit(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, EOK)
}

func ReadExit(r io.Reader) (exitCode ExitCode, errMsg string, err error) {
	err = binary.Read(r, binary.LittleEndian, &exitCode)
	if err != nil {
		return
	}
	if exitCode != EOK {
		errMsg, err = ReadString(r)
	}
	return
}

func WriteString(w io.Writer, str string) error {
	length := int64(len(str))
	err := binary.Write(w, binary.LittleEndian, length)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, str)
	return err
}

func ReadString(r io.Reader) (string, error) {
	var length int64
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	_, err = io.CopyN(&sb, r, length)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
