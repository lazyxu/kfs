package common

type Client[T any] interface {
	Server() *EventServer[T]
	Chan() chan T
	Message() T
}

type DefaultClient[T any] struct {
	s  *EventServer[T]
	ch chan T
}

func (c *DefaultClient[T]) Server() *EventServer[T] {
	return c.s
}

func (c *DefaultClient[T]) Chan() chan T {
	return c.ch
}

func (c *DefaultClient[T]) Message() T {
	return c.s.Message()
}
