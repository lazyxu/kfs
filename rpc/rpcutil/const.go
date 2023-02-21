package rpcutil

type CommandType = int8

const (
	CommandPing CommandType = iota

	CommandReset

	CommandList
	CommandOpen

	CommandUpload
	CommandUploadDirItem
	CommandTouch

	CommandDownload

	CommandCat

	CommandBranchCheckout
	CommandBranchInfo
	CommandBranchList
)

type Status = int8

const (
	EOK Status = iota
	EInvalid
)
