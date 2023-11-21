package core

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"sync"
)

type EventClient[T any] interface {
	Chan() chan T
	Message() T
}

type EventServer[T any] struct {
	Clients   sync.Map // map[*http.Request]EventClient
	NewClient func(c echo.Context, kfsCore *KFS) (EventClient[T], error)
}

func (s *EventServer[T]) Add(c echo.Context, kfsCore *KFS) (EventClient[T], error) {
	client, err := s.NewClient(c, kfsCore)
	if err != nil {
		return nil, err
	}
	s.Clients.Store(c, client)
	return client, nil
}

func (s *EventServer[T]) Delete(c echo.Context) {
	client, _ := s.Clients.Load(c)
	// Should be deleted before closing.
	s.Clients.Delete(c)
	close(client.(EventClient[T]).Chan())
}

func (s *EventServer[T]) Handle(c echo.Context, kfsCore *KFS) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	client, err := s.Add(c, kfsCore)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	defer func() {
		s.Delete(c)
		fmt.Println("Closing connection")
	}()

	msg := client.Message()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("invalid json msg: %+v\n", data)
	}
	fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
	c.Response().Flush()

	for {
		select {
		case msg = <-client.Chan():
			data, err = json.Marshal(msg)
			if err != nil {
				log.Panicf("invalid msg: %+v\n", msg)
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
			c.Response().Flush()

		case <-c.Request().Context().Done():
			fmt.Println("Connection closed")
			return nil
		}
	}
}

func (s *EventServer[T]) SendAll() {
	s.Clients.Range(func(key, value any) bool {
		client := value.(EventClient[T])
		// TODO: may panic here.
		client.Chan() <- client.Message()
		return true
	})
}
