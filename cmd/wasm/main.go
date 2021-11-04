package wasm

import (
	"syscall/js"
)

func add(i []js.Value) {
	valueA := i[0].Int()
	valueB := i[1].Int()
	sum := valueA + valueB
	js.Global().Get("console").Call("log", "sum: ", js.ValueOf(sum))
}

func main() {
	c := make(chan struct{}, 0)
	callback := js.NewCallback(add)
	js.Global().Set("add", callback)
	defer callback.Release()
	<-c
}
