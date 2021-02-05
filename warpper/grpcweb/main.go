package grpcweb

import (
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/lazyxu/kfs/warpper/grpcweb/pb"
	"github.com/lazyxu/kfs/warpper/grpcweb/rootdirectory"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(httpPort int, s storage.Storage) {
	logrus.SetLevel(logrus.TraceLevel)
	Init(s)
	go rootdirectory.Socket(s, 9877)
	fsServer := rootdirectory.New(s)
	serverHttp(getHandler(fsServer, s), ":"+strconv.Itoa(httpPort))
}

func getHandler(fsServer pb.KoalaFSServer, s storage.Storage) http.Handler {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/api/ref/:name", func(c echo.Context) error {
		hash, err := s.GetRef(c.Param("name"))
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, hash)
	})
	e.PUT("/api/ref/:name/:old/:new", func(c echo.Context) error {
		err := s.UpdateRef(c.Param("name"), c.Param("old"), c.Param("new"))
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "OK")
	})
	e.POST("/api/upload", func(c echo.Context) error {
		hash, err := obj.WriteBlob(c.Request().Body)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, hash)
	})
	e.GET("/api/download/:hash", func(c echo.Context) error {
		hash := c.Param("hash")
		err := obj.ReadBlob(hash, func(r io.Reader) error {
			_, err := io.Copy(c.Response(), r)
			return err
		})
		if err != nil {
			return err
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})
	mux := http.NewServeMux()
	mux.Handle("/", e)
	server := grpc.NewServer()
	pb.RegisterKoalaFSServer(server, fsServer)
	go func() {
		server := grpc.NewServer()
		pb.RegisterKoalaFSServer(server, fsServer)
		lis, err := net.Listen("tcp", ":9092")
		if err != nil {
			logrus.Fatal("failed to listen", err)
			return
		}
		if err := server.Serve(lis); err != nil {
			logrus.Fatal("failed to serve", err)
		}
	}()
	wrappedGrpc := grpcweb.WrapServer(server)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "kfs-mount", "X-gRPC-Web", "X-CSRF-Token"},
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
