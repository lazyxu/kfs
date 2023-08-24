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
	EOK       Status = 0
	ENotExist Status = 1
	EInvalid  Status = 0xF
)
