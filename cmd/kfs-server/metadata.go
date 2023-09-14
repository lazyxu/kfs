package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/rpc/server"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Client struct {
	sseChannel     chan string
	sseJsonChannel chan interface{}
}

var taskListClients sync.Map // map[*http.Request]*Client

func apiEventMetadataAnalysisTask(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	sseChannel := make(chan string)
	sseJsonChannel := make(chan interface{})
	taskListClients.Store(c, &Client{
		sseChannel:     sseChannel,
		sseJsonChannel: sseJsonChannel,
	})

	defer func() {
		close(sseChannel)
		close(sseJsonChannel)
		taskListClients.Delete(c)
		fmt.Println("Closing connection")
	}()

	data, err := json.Marshal(metadataAnalysisTask)
	if err != nil {
		log.Panicf("invalid obj: %+v\n", metadataAnalysisTask)
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

var (
	StatusIdle           = 0
	StatusFinished       = 1
	StatusCanceled       = 2
	StatusError          = 3
	StatusWaitRunning    = 4
	StatusWaitCanceled   = 5
	StatusRunningCollect = 6
	StatusRunningAnalyze = 7
)

type MetadataAnalysisTask struct {
	cancel context.CancelFunc
	Status int      `json:"status"`
	Cnt    int      `json:"cnt"`
	Total  int      `json:"total"`
	Errors []string `json:"errors"`
}

var metadataAnalysisTask = MetadataAnalysisTask{
	cancel: nil,
	Status: StatusIdle,
	Cnt:    0,
	Total:  0,
	Errors: make([]string, 0),
}

var metadataAnalysisTaskMutex = &sync.RWMutex{}

func noteTaskListToClients() {
	taskListClients.Range(func(key, value any) bool {
		client := value.(*Client)
		client.sseJsonChannel <- metadataAnalysisTask
		return true
	})
}

func setTaskStatus(status int) {
	metadataAnalysisTaskMutex.Lock()
	metadataAnalysisTask.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		metadataAnalysisTask.cancel = nil
	}
	if status == StatusWaitRunning {
		metadataAnalysisTask.Errors = make([]string, 0)
		metadataAnalysisTask.Cnt = 0
		metadataAnalysisTask.Total = 0
	}
	metadataAnalysisTaskMutex.Unlock()
	noteTaskListToClients()
}

func setTaskTotal(total int) {
	metadataAnalysisTaskMutex.Lock()
	metadataAnalysisTask.Status = StatusRunningAnalyze
	metadataAnalysisTask.Total = total
	metadataAnalysisTaskMutex.Unlock()
	noteTaskListToClients()
}

func addTaskCnt() {
	metadataAnalysisTaskMutex.Lock()
	metadataAnalysisTask.Cnt++
	metadataAnalysisTaskMutex.Unlock()
	noteTaskListToClients()
}

func addTaskError(err error) {
	metadataAnalysisTaskMutex.Lock()
	metadataAnalysisTask.Errors = append(metadataAnalysisTask.Errors, err.Error())
	metadataAnalysisTaskMutex.Unlock()
	noteTaskListToClients()
}

func apiStartMetadataAnalysisTask(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	metadataAnalysisTaskMutex.Lock()
	defer metadataAnalysisTaskMutex.Unlock()
	if !start {
		setTaskStatus(StatusWaitCanceled)
		metadataAnalysisTask.cancel()
		return c.String(http.StatusOK, "")
	}
	tryStartMetadataAnalysisTask()
	return c.String(http.StatusOK, "")
}

func tryStartMetadataAnalysisTask() {
	if metadataAnalysisTask.Status == StatusWaitRunning ||
		metadataAnalysisTask.Status == StatusRunningCollect ||
		metadataAnalysisTask.Status == StatusRunningAnalyze ||
		metadataAnalysisTask.Status == StatusWaitCanceled {
		return
	}
	metadataAnalysisTask.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	metadataAnalysisTask.cancel = cancel
	go func() {
		err := analyzeMetadata(ctx)
		if err == nil {
			setTaskStatus(StatusFinished)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			setTaskStatus(StatusCanceled)
			return
		}
		setTaskStatus(StatusError)
	}()
}

func analyzeMetadata(ctx context.Context) error {
	setTaskStatus(StatusRunningCollect)
	hashList, err := kfsCore.Db.ListExpectFileType(ctx)
	if err != nil {
		return err
	}
	setTaskTotal(len(hashList))
	for _, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		ft, err := server.AnalyzeFileType(kfsCore, hash)
		if err != nil {
			addTaskError(err)
			continue
		}
		err = server.InsertExif(context.TODO(), kfsCore, hash, ft)
		if err != nil {
			addTaskError(err)
			continue
		}
		err = server.InsertFileType(context.TODO(), kfsCore, hash, ft)
		if err != nil {
			addTaskError(err)
			continue
		}
		addTaskCnt()
	}
	if err != nil {
		return err
	}
	return nil
}
