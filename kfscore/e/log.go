package e

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// fnName returns the name of the calling +2 function
func fnName() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "*Unknown*"
	}
	filename := path.Base(file)
	name := runtime.FuncForPC(pc).Name()
	dot := strings.LastIndex(name, ".")
	if dot >= 0 {
		if name[dot-1] == ')' {
			dot2 := strings.LastIndex(name, "(")
			if dot2 >= 0 {
				dot = dot2 - 1
			}
		}
		name = name[dot+1:]
	}
	return "[" + filename + ":" + strconv.Itoa(line) + "] " + name
}

// Trace debugs the entry and exit of the calling function
//
// It is designed to be used in a defer statement so it returns a
// function that logs the exit parameters.
//
// Any pointers in the exit function will be dereferenced
func Trace(enterFields logrus.Fields) func(exitFields func() logrus.Fields) {
	name := fnName()
	logrus.WithFields(enterFields).Trace(name)
	return func(exitFields func() logrus.Fields) {
		if exitFields == nil {
			return
		}
		fields := exitFields()
		if fields == nil {
			return
		}
		logrus.WithFields(fields).Trace(name)
	}
}
