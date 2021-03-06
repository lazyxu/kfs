package grpcweb

import (
	"crypto/tls"
	"log"
	"net/http"
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
