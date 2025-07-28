package local_file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/db/sqlc"
	"github.com/robfig/cron/v3"
)

func jsonBody(m any) io.Reader {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

var DeviceId string

func NewDeviceIfNeeded(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	socketServer := json.ConfigCompatibleWithStandardLibrary.Get(data, "socketServer").ToString()
	if socketServer == "" {
		fmt.Printf("请配置socketServer\n")
		os.Exit(1)
	}
	webServer := json.ConfigCompatibleWithStandardLibrary.Get(data, "webServer").ToString()
	if webServer == "" {
		fmt.Printf("请配置webServer\n")
		os.Exit(1)
	}
	id := json.ConfigCompatibleWithStandardLibrary.Get(data, "deviceId").ToString()
	if id != "" {
		DeviceId = id
	} else {
		id, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		DeviceId = id.String()
		m := map[string]interface{}{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			panic(err)
		}
		m["deviceId"] = DeviceId
		data, err = json.MarshalIndent(m, "", " ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(configPath, data, 0o600)
		if err != nil {
			panic(err)
		}
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	os := runtime.GOOS
	name := hostname + "@" + os
	buf := jsonBody(map[string]interface{}{
		"id":       DeviceId,
		"name":     name,
		"os":       os,
		"hostname": hostname,
	})
	_, err = http.Post(webServer+"/api/v1/devices", "application/json", buf)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(webServer + "/api/v1/listLocalFileDriver?deviceId=" + DeviceId)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	code := json.ConfigCompatibleWithStandardLibrary.Get(body, "code").ToInt()
	if code != 0 {
		panic("error code: " + strconv.Itoa(code))
	}
	var list []sqlc.LocalFileSyncToLogicFile
	json.ConfigCompatibleWithStandardLibrary.Get(body, "data").ToVal(&list)
	for _, d := range list {
		if d.Enabled && d.DeviceID == DeviceId && d.SyncAtStart {
			StartLocalFileSync(d.ID, socketServer, d.Cron, d.Srcpath, d.Ignores, d.Encoder)
		}
	}
	return DeviceId
}

var cronTasks sync.Map

type CronTask struct {
	c      *cron.Cron
	id     cron.EntryID
	cancel context.CancelFunc
}

func StartLocalFileSync(fileID int64, serverAddr string, cronStr string, srcPath, ignores, encoder string) {
	actual, loaded := cronTasks.LoadOrStore(fileID, &CronTask{
		c:      cron.New(),
		id:     -1,
		cancel: nil,
	})
	t := actual.(*CronTask)
	if loaded {
		if t.cancel != nil {
			t.cancel()
			t.cancel = nil
		}
		t.c.Remove(t.id)
	}
	var err error
	t.id, err = t.c.AddFunc(cronStr, func() {
		ctx, cancel := context.WithCancel(context.TODO())
		t.cancel = cancel
		d, err := GetOrLoad(fileID)
		if err != nil {
			cronTasks.LoadAndDelete(fileID)
			return
		}
		d.StartOrStop(ctx, true, serverAddr, srcPath, ignores, encoder)
	})
	if err != nil {
		panic(err)
	}
	t.c.Start()
}
