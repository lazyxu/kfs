protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. cmd/server/pb/fs.proto
cp cmd/server/pb/*.pb.go cmd/client/pb
