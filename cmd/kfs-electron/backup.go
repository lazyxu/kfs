package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func apiBackupTask(c echo.Context) error {
	list, err := db.ListBackupTask(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, list)
}

type Client struct {
	sseChannel     chan string
	sseJsonChannel chan interface{}
}

var clients sync.Map // map[*http.Request]*Client

func apiEventBackupTask(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	sseChannel := make(chan string)
	sseJsonChannel := make(chan interface{})
	client := &Client{
		sseChannel:     sseChannel,
		sseJsonChannel: sseJsonChannel,
	}
	clients.Store(c, client)

	defer func() {
		close(sseChannel)
		close(sseJsonChannel)
		clients.Delete(c)
		fmt.Println("Closing connection")
	}()

	obj, err := getTaskList(c.Request().Context())
	if err != nil {
		return err
	}
	data, err := json.Marshal(obj)
	if err != nil {
		log.Panicf("invalid obj: %+v\n", obj)
	}
	fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
	c.Response().Flush()

	for {
		select {
		case msg := <-sseChannel:
			fmt.Fprintf(c.Response(), "data: %s\n\n", msg)
			c.Response().Flush()

		case obj := <-sseJsonChannel:
			data, err := json.Marshal(obj)
			if err != nil {
				log.Panicf("invalid obj: %+v\n", obj)
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
			c.Response().Flush()

		case <-c.Request().Context().Done():
			fmt.Println("Connection closed")
			return nil
		}
	}
}

func getTaskList(ctc context.Context) (list []BackupTask, err error) {
	list, err = db.ListBackupTask(ctc)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func noteTaskListToClients() {
	clients.Range(func(key, value any) bool {
		c := key.(echo.Context)
		client := value.(*Client)
		list, err := getTaskList(c.Request().Context())
		if err != nil {
			c.Logger().Error(err)
			return true
		}
		client.sseJsonChannel <- list
		return true
	})
}

func apiNewBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	srcPath := c.QueryParam("srcPath")
	driverName := c.QueryParam("driverName")
	dstPath := c.QueryParam("dstPath")
	encoder := c.QueryParam("encoder")
	concurrentStr := c.QueryParam("concurrent")
	concurrent, err := strconv.Atoi(concurrentStr)
	if err != nil {
		return err
	}
	err = upsertBackup(c.Request().Context(), db, name, description, srcPath, driverName, dstPath, encoder, concurrent)
	if err != nil {
		return err
	}
	noteTaskListToClients()
	return c.String(http.StatusOK, "")
}

func apiDeleteBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	err := db.DeleteBackupTask(c.Request().Context(), name)
	if err != nil {
		return err
	}
	noteTaskListToClients()
	return c.String(http.StatusOK, "")
}
