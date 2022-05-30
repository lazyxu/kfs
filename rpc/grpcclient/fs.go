package grpcclient

type GRPCFS struct {
	RemoteAddr string
}

func New(remoteAddr string) *GRPCFS {
	return &GRPCFS{
		RemoteAddr: remoteAddr,
	}
}
