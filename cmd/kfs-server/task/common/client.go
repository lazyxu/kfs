package common

type Client[T any] interface {
	Chan() chan T
	Message() T
}
