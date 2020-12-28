package grpcweb

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

func serverHttp(handler http.Handler, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatal("failed to listen", err)
	}
	httpServer := &http.Server{
		Addr:    httpsPort,
		Handler: http.DefaultServeMux,
	}
	httpServer.Handler = handler
	logrus.WithFields(logrus.Fields{"port": port}).Info("Listening http")
	if err := httpServer.Serve(lis); err != nil {
		logrus.Fatal("failed to serve", err)
	}
}
