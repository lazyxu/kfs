package client

type GRPCFS struct {
	RemoteAddr string
}

func New(remoteAddr string) *GRPCFS {
	return &GRPCFS{
		RemoteAddr: remoteAddr,
	}
}

func (fs *GRPCFS) Close() error {
	return nil
}
