package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Backup(server pb.KoalaFS_BackupServer) (err error) {
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return
	}
	defer kfsCore.Close()
	req := &pb.BackupReq{}
	var h *pb.BackupReqHeader
	var m *pb.UploadReqMetadata
	var exist bool
	fmt.Println("-----------")
	for {
		req, err = server.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		if req.Done {
			break
		}
		fmt.Println("Backup", req.Header, req.Metadata)
		if req.Header != nil {
			h = req.Header
		}
		if req.Metadata == nil {
			continue
		}
		m = req.Metadata
		if os.FileMode(m.Mode).IsRegular() {
			exist, err = kfsCore.S.WriteFn(m.Hash, func(f io.Writer, hasher io.Writer) error {
				for {
					fmt.Println("server.Recv")
					req, err = server.Recv()
					if req != nil {
						fmt.Println("Backup", h.Base, req.IsLast, len(req.Bytes))
					}
					if err != nil && err != io.EOF {
						return err
					}
					if err == io.EOF {
						return nil
					}
					_, err = hasher.Write(req.Bytes)
					if err != nil {
						return nil
					}
					_, err = f.Write(req.Bytes)
					if err != nil {
						return nil
					}
					if req.IsLast {
						return nil
					}
				}
			})
			if err != nil {
				return
			}
			fmt.Println("Backup", m, exist)
			err = server.Send(&pb.BackupResp{
				UploadResp: &pb.UploadResp{
					Exist: exist,
				},
			})
			if err != nil {
				return
			}
		} else {
			// TODO: upload dir
		}
	}
	fmt.Println("Backup finish")
	err = server.Send(&pb.BackupResp{
		Done:    true,
		Err:     nil,
		Process: "",
	})
	return
}
