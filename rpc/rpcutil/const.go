package rpcutil

type CommandType = int8

const (
	CommandPing CommandType = iota

	CommandReset

	CommandList

	CommandUpload
	CommandTouch

	CommandDownload

	CommandCat

	CommandBranchCheckout
	CommandBranchInfo
)

type Status = int8

const (
	EOK Status = iota
	EInvalid
)
