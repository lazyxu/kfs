package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/labstack/gommon/log"

	"github.com/lazyxu/kfs/cmd/client/localfs"

	"github.com/gorilla/websocket"
	"github.com/thedevsaddam/gojsonq/v2"
)

type JsonQ struct {
	*gojsonq.JSONQ
}

func NewJsonQ(data string) *JsonQ {
	return &JsonQ{gojsonq.New().FromString(data)}
}

func (j *JsonQ) RFindStringOrDefault(path string, defaultVal string) string {
	j.Reset()
	i := j.Find(path)
	if s, ok := i.(string); ok {
		return s
	}
	return defaultVal
}

func (j *JsonQ) RFindString(path string) (string, error) {
	j.Reset()
	i := j.Find(path)
	if s, ok := i.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("%s should be a string", path)
}

type MethodFunc func(ctx context.Context, q *JsonQ, send func(v interface{}) error) error

var methods map[string]MethodFunc

func init() {
	methods = make(map[string]MethodFunc)
}

func RegisterInvokeMethod(method string, fn func(q *JsonQ) (interface{}, error)) {
	methods[method] = func(ctx context.Context, q *JsonQ, send func(v interface{}) error) error {
		id, err := q.RFindString("id")
		if err != nil {
			return err
		}
		data, err := fn(q)
		if err != nil {
			return err
		}
		return send(localfs.Result{
			ID:     id,
			Result: data,
		})
	}
}

func Register1ton(method string, fn func(ctx context.Context, q *JsonQ, ch chan<- interface{}) error) {
	methods[method] = func(ctx context.Context, q *JsonQ, send func(v interface{}) error) error {
		id, err := q.RFindString("id")
		if err != nil {
			return err
		}
		ch := make(chan interface{})
		go func() {
			for {
				select {
				case <-ctx.Done():
					logrus.Debug("context deadline exceed!")
					return
				case data := <-ch:
					err = send(localfs.Result{
						ID:     id,
						Result: data,
					})
					if err != nil {
						log.Error(err)
					}
				}
			}
		}()
		err = fn(ctx, q, ch)
		if err != nil {
			return err
		}
		return nil
	}
}

func myJsonHandler(ctx context.Context, data string, send func(v interface{}) error) error {
	q := NewJsonQ(data)
	method, err := q.RFindString("method")
	if err != nil {
		return err
	}
	fmt.Printf("ReadMessage: %s\n", data)
	if fn, ok := methods[method]; ok {
		return fn(ctx, q, send)
	}
	switch method {
	case "echo":
		return send(data)
	default:
		fmt.Println("Invalid method:", method)
	}
	return nil
}

func JsonHandleWrapper(jsonHandler func(ctx context.Context, data string, send func(i interface{}) error) error) Handle {
	return func(ctx context.Context, messageType int, data []byte, send func(messageType int, data []byte) error) error {
		if messageType == websocket.TextMessage && data != nil {
			return jsonHandler(ctx, string(data), func(v interface{}) error {
				if s, ok := v.(string); ok {
					return send(websocket.TextMessage, []byte(s))
				}
				data, err := json.Marshal(v)
				if err != nil {
					return err
				}
				return send(websocket.TextMessage, data)
			})
		}
		return nil
	}
}

func Start() {
	startWS(JsonHandleWrapper(myJsonHandler))
}
