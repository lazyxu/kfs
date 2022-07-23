package main

import "errors"

const (
	ConfigFileStr   = "config-file"
	ServerAddrStr   = "server-addr"
	SocketServerStr = "socket-server"
	GrpcServerStr   = "grpc-server"

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
