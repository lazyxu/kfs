package rpcutil

type CommandType = int8

const (
	CommandPing CommandType = iota

	CommandUpload

	CommandDownload
)

type ExitCode = int8

const (
	EOK ExitCode = iota
	EInvalid
)
