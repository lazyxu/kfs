package kfsserver

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func catch(err *error) {
	if *err != nil {
		logrus.Errorln("Result", (*err).Error())
	}
	if val := recover(); val != nil {
		var errMsg string
		var stack string
		if e, ok := val.(error); ok {
			errMsg = e.Error()
			stack = trimStack("/src/runtime/panic.go")
			*err = fmt.Errorf("internal Server Error")
		} else if e, ok := val.(*logrus.Entry); ok {
			errMsg = e.Message
			stack = trimStack("/exported.go")
			*err = fmt.Errorf("internal Server Error")
		} else {
			*err = fmt.Errorf("internal Server Error")
			errMsg = fmt.Sprintf("unknown error: %s %v", reflect.TypeOf(val), val)
		}
		logrus.Errorln("Panic", errMsg+"\n"+stack)
	}
}

func trimStack(s string) string {
	e := []byte("pb.go")
	line := []byte("\n")
	stack := debug.Stack()
	length := len(stack)
	start := bytes.Index(stack, []byte(s))
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}
