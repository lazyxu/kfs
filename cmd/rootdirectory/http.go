package main

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

func serverHttp(handler http.Handler) {
	lis, err := net.Listen("tcp", httpPort)
	if err != nil {
		logrus.Fatal("failed to listen", err)
	}
	httpServer := &http.Server{
		Addr:    httpsPort,
		Handler: http.DefaultServeMux,
	}
	httpServer.Handler = handler
	logrus.WithFields(logrus.Fields{"httpPort": httpPort}).Info("Listening")
	if err := httpServer.Serve(lis); err != nil {
		logrus.Fatal("failed to serve", err)
	}
}
