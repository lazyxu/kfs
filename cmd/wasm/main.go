package wasm

import (
	"strings"
	"syscall/js"

	"github.com/lazyxu/kfs/kfscore/storage"
)

func add(this js.Value, args []js.Value) interface{} {
	println("wasm: add")
	d := &storage.Directory{}
	storage.DefaultDirectoryEncoderDecoder.Decode(d, strings.NewReader(""))
	valueA := args[0].Int()
	valueB := args[1].Int()
	sum := valueA + valueB
	js.Global().Get("console").Call("log", "sum: ", js.ValueOf(sum))
	return nil
}

func main() {
	println("load go-wasm module")
	c := make(chan struct{}, 0)
	callback := js.FuncOf(add)
	js.Global().Set("add", callback)
	defer callback.Release()
	<-c
}
