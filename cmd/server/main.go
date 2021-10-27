package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/lazyxu/kfs/cmd/server/kfsserver"

	"github.com/lazyxu/kfs/cmd/server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	httpsPort = ":9091"
)

func serverHttps(handler http.Handler) {
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
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	tlsHttpServer.Handler = handler
	log.Fatal(tlsHttpServer.ListenAndServeTLS("cert.pem", "key.pem"))
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	var opts []grpc.ServerOption
	fsServer := kfsserver.New()
	server := grpc.NewServer(opts...)
	pb.RegisterKoalaFSServer(server, fsServer)
	logrus.Println("listening on port ", 9092)
	lis, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		logrus.Fatal("failed to listen", err)
		return
	}
	if err := server.Serve(lis); err != nil {
		logrus.Fatal("failed to serve", err)
	}
}
