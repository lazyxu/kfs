package utils

import "errors"

const (
	ServerTypeStr    = "server-type"
	ServerTypeLocal  = "local"
	ServerTypeRemote = "remote"
	ServerAddrStr    = "server-addr"

	HumanizeStr    = "humanize"
	BackupPathStr  = "backup-path"
	BranchNameStr  = "branch-name"
	PathStr        = "path"
	ChunkSizeStr   = "block-bytes"
	DescriptionStr = "description"
)

var (
	InvalidServerType = errors.New("invalid server type")
)
