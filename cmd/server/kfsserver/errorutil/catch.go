package errorutil

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
)

var SystemInternalError = status.Errorf(codes.Internal, "系统内部错误")

func Catch(err *error) {
	if *err != nil {
		logrus.Errorln("Result", (*err).Error())
	}
	if val := recover(); val != nil {
		var errMsg string
		var stack string
		if e, ok := val.(error); ok {
			errMsg = e.Error()
			stack = trimStack("/src/runtime/panic.go")
			*err = SystemInternalError
		} else if e, ok := val.(*logrus.Entry); ok {
			errMsg = e.Message
			stack = trimStack("/exported.go")
			*err = SystemInternalError
		} else {
			*err = SystemInternalError
			errMsg = fmt.Sprintf("unknown error: %s %v", reflect.TypeOf(val), val)
		}
		logrus.Error(errMsg + "\n" + stack)
	}
}

func Max(x int, y int) (max int) {
	if x > y {
		return x
	}
	return y
}

func trimStack(s string) string {
	line := []byte("\n")
	stack := debug.Stack()
	start1 := bytes.Index(stack, []byte("errorutil/error.go"))
	start2 := bytes.Index(stack, []byte(s))
	stack = stack[Max(start1, start2):]
	end := bytes.Index(stack, []byte("pb.go"))
	if end != -1 {
		stack = stack[:end]
		end = bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
		end = bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
	}
	return string(stack)
}
