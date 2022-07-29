package rpcutil

type CommandType = int8

const (
	CommandPing CommandType = iota

	CommandList

	CommandUpload
	CommandTouch

	CommandDownload

	CommandCat
)

type Status = int8

const (
	EOK Status = iota
	EInvalid
)
