package main

import "errors"

const (
	ConfigFileStr   = "config-file"
	SocketServerStr = "socket-server"

	HumanizeStr    = "humanize"
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
