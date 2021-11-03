package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/cmd/client/pb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lazyxu/kfs/cmd/client/kfsclient"
)

type Config struct {
	ClientID string
}

func GetConfig() (*Config, error) {
	file, err := ioutil.ReadFile("kfs-config.json")
	if err != nil {
		return nil, err
	}
	data := &Config{}
	err = json.Unmarshal(file, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

var gClient *kfsclient.Client
var once sync.Once

func GetClient() *kfsclient.Client {
	once.Do(func() {
		gClient = kfsclient.New("localhost:9092")
	})
	return gClient
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, this is kfs client!")
	})
	e.GET("/api/clientID", func(c echo.Context) error {
		config, err := GetConfig()
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "读取配置失败")
		}
		return c.String(http.StatusOK, config.ClientID)
	})
	e.POST("/api/getBranchHash", func(c echo.Context) error {
		branch := &pb.Branch{}
		err := c.Bind(branch)
		if err != nil {
			return err
		}
		hash, err := GetClient().PbClient.GetBranchHash(context.TODO(), branch)
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.JSON(http.StatusOK, hash.Hash)
	})
	e.POST("/api/listBranches", func(c echo.Context) error {
		branches, err := GetClient().ListBranches(context.TODO())
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.JSON(http.StatusOK, branches)
	})
	e.POST("/api/createBranch", func(c echo.Context) error {
		branch := &pb.Branch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.CreateBranch(context.TODO(), branch)
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.NoContent(http.StatusOK)
	})
	e.POST("/api/deleteBranch", func(c echo.Context) error {
		branch := &pb.Branch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.DeleteBranch(context.TODO(), branch)
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.NoContent(http.StatusOK)
	})
	e.POST("/api/renameBranch", func(c echo.Context) error {
		branch := &pb.RenameBranch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.RenameBranch(context.TODO(), branch)
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.NoContent(http.StatusOK)
	})
	e.POST("/api/readObject", func(c echo.Context) error {
		req := &pb.Hash{}
		err := c.Bind(req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		client, err := GetClient().PbClient.ReadObject(context.TODO(), req)
		if err != nil {
			return FromGrpcError(c, err)
		}
		for {
			chunk, err := client.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return FromGrpcError(c, err)
			}
			_, err = c.Response().Write(chunk.GetChunk())
			if err != nil {
				return err
			}
		}
		c.Response().Flush()
		return nil
	})
	e.GET("/api/connect", func(c echo.Context) error {
		fmt.Println("connect")
		client := GetClient()
		hash, err := client.WriteObject(context.TODO(), []byte("111"))
		if err != nil {
			return FromGrpcError(c, err)
		}
		err = client.ReadObject(context.TODO(), hex.EncodeToString(hash), func(buf []byte) error {
			fmt.Println("ReadObject", string(buf))
			return nil
		})
		if err != nil {
			return FromGrpcError(c, err)
		}
		return c.NoContent(http.StatusOK)
	})
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	e.Logger.Fatal(e.StartTLS(":"+port, "localhost.pem", "localhost-key.pem"))
}
