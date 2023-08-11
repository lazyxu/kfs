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

	CommandUploadV2File
	CommandUploadV2Dir

	CommandDownload

	CommandCat

	CommandBranchCheckout
	CommandBranchInfo
	CommandBranchList
)

type Status = int8

const (
	EOK Status = iota
	ENotExist
	EInvalid = -1
)
