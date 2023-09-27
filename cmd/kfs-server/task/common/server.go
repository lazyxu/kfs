package common

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"sync"
)

type EventServer[T any] struct {
	Clients sync.Map // map[*http.Request]Client
	Message func() T
}

func (s *EventServer[T]) Add(c echo.Context) Client[T] {
	client := &DefaultClient[T]{
		ch: make(chan T),
	}
	s.Clients.Store(c, client)
	return client
}

func (s *EventServer[T]) Delete(c echo.Context) {
	client, _ := s.Clients.Load(c)
	close(client.(Client[T]).Chan())
	s.Clients.Delete(c)
}

func (s *EventServer[T]) Handle(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	client := s.Add(c)

	defer func() {
		s.Delete(c)
		fmt.Println("Closing connection")
	}()

	msg := client.Message()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("invalid msg: %+v\n", data)
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
		client := value.(Client[T])
		client.Chan() <- s.Message()
		return true
	})
}
