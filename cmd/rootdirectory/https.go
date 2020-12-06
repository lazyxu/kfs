package main

import (
	"crypto/tls"
	"log"
	"math"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/lazyxu/kfs/pb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

const (
	httpPort  = ":9090"
	httpsPort = ":9091"
)

func serverHttps(srv pb.KoalaFSServer) {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example pkg.\n"))
	})
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	tlsHttpServer := &http.Server{
		Addr:         httpsPort,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(math.MaxInt32))
	pb.RegisterKoalaFSServer(grpcServer, srv)
	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	tlsHttpServer.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsAcceptableGrpcCorsRequest(req) || wrappedGrpc.IsGrpcWebRequest(req) {
			cors.New(cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "kfs-pwd", "kfs-mount", "X-gRPC-Web", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300, // Maximum value not ignored by any of major browsers
				Debug:            false,
			}).HandlerFunc(resp, req)
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		log.Println("cors")
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}).HandlerFunc(resp, req)

		// Fall back to other servers.
		http.DefaultServeMux.ServeHTTP(resp, req)
	})
	log.Fatal(tlsHttpServer.ListenAndServeTLS("cert.pem", "key.pem"))
}
