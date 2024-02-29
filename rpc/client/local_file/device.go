package local_file

import (
	"bytes"
	"context"
	"fmt"
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

func NewDeviceIfNeeded(configPath string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	var deviceId string
	if id, o := m["deviceId"]; o {
		deviceId = id.(string)
	} else {
		deviceId = ""
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	os := runtime.GOOS
	name := hostname + "@" + os
	webServerHost := m["webServer"].(string)
	buf := jsonBody(map[string]interface{}{
		"id":       deviceId,
		"name":     name,
		"os":       os,
		"hostname": hostname,
	})
	_, err = http.Post(webServerHost+"/api/v1/devices", "application/json", buf)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(webServerHost + "/api/v1/listLocalFileDriver?deviceId=" + deviceId)
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
		StartLocalFileSync(d.Id, m["socketServer"].(string), d.H, d.M, d.SrcPath, d.Ignores, d.Encoder)
	}
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
