package local_file

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/dao"
	"github.com/robfig/cron/v3"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func jsonBody(m any) io.Reader {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

var deviceId string

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
		deviceId = id
	} else {
		id, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		deviceId = id.String()
		m := map[string]interface{}{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			panic(err)
		}
		m["deviceId"] = deviceId
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
		"id":       deviceId,
		"name":     name,
		"os":       os,
		"hostname": hostname,
	})
	_, err = http.Post(webServer+"/api/v1/devices", "application/json", buf)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(webServer + "/api/v1/listLocalFileDriver?deviceId=" + deviceId)
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
	var drivers []dao.Driver
	json.ConfigCompatibleWithStandardLibrary.Get(body, "data").ToVal(&drivers)
	for _, d := range drivers {
		StartLocalFileSync(d.Id, socketServer, d.H, d.M, d.SrcPath, d.Ignores, d.Encoder)
	}
	return deviceId
}

var cronTasks sync.Map

type CronTask struct {
	c      *cron.Cron
	id     cron.EntryID
	cancel context.CancelFunc
}

func StartLocalFileSync(driverId uint64, serverAddr string, h int64, m int64, srcPath, ignores, encoder string) {
	actual, loaded := cronTasks.LoadOrStore(driverId, &CronTask{
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
	spec := fmt.Sprintf("%d %d * * ?", m, h)
	var err error
	t.id, err = t.c.AddFunc(spec, func() {
		ctx, cancel := context.WithCancel(context.TODO())
		t.cancel = cancel
		d, err := GetOrLoadDriver(driverId)
		if err != nil {
			cronTasks.LoadAndDelete(driverId)
			return
		}
		d.StartOrStop(ctx, true, serverAddr, srcPath, ignores, encoder)
	})
	if err != nil {
		panic(err)
	}
	t.c.Start()
}
