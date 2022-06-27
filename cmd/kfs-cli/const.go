package main

import "errors"

const (
	ServerTypeStr       = "server-type"
	ServerTypeLocal     = "local"
	ServerTypeRemote    = "remote"
	ServerAddrStr       = "server-addr"
	GrpcServerAddrStr   = "grpc-server-addr"
	SocketServerAddrStr = "socket-server-addr"

	HumanizeStr    = "humanize"
	BackupPathStr  = "backup-path"
	BranchNameStr  = "branch-name"
	PathStr        = "path"
	DirPathStr     = "dir-path"
	ChunkSizeStr   = "block-bytes"
	DescriptionStr = "description"

	VerboseStr     = "verbose"
	ConcurrentStr  = "concurrent"
	EncoderStr     = "encoder"
	CpuProfilerStr = "cpu-profile"
)

var (
	InvalidServerType = errors.New("invalid server type")
)
