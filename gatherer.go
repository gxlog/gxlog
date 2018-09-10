package gxlog

import (
	"runtime"
	"time"
)

type gatherer struct{}

func (this *gatherer) gather(calldepth int, level LogLevel, msg string) *Record {
	now := time.Now()
	var funcName string
	pc, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "?"
		line = -1
		funcName = "?"
	} else {
		funcName = runtime.FuncForPC(pc).Name()
	}
	return &Record{
		Time:     now,
		Level:    level,
		Pathname: file,
		Line:     line,
		Func:     funcName,
		Msg:      msg,
	}
}
