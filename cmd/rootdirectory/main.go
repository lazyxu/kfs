package main

import (
	"io"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rootdirectory"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	Init()
	srv := rootdirectory.New(s)
	serverHttp(getHandler(srv))
}

func getHandler(srv pb.KoalaFSServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example pkg.\n"))
	})
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, req *http.Request) {
		hash, err := obj.WriteBlob(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hash))
	})
	mux.HandleFunc("/api/download", func(w http.ResponseWriter, req *http.Request) {
		hash := req.URL.Query().Get("hash")
		err := obj.ReadBlob(hash, func(r io.Reader) error {
			_, err := io.Copy(w, r)
			return err
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, srv)
	wrappedGrpc := grpcweb.WrapServer(s)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "kfs-pwd", "kfs-mount", "X-gRPC-Web", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
			Debug:            false,
		}).HandlerFunc(resp, req)
		if wrappedGrpc.IsAcceptableGrpcCorsRequest(req) || wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		mux.ServeHTTP(resp, req)
	})
}
