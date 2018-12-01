package gxlog

import (
	"runtime"
	"time"
)

type gatherer struct {
	record Record
}

func (this *gatherer) Gather(calldepth int, level LogLevel, msg string) *Record {
	this.record.Time = time.Now()
	this.record.Level = level
	this.record.Msg = msg

	var funcName string
	pc, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "?"
		line = -1
		funcName = "?"
	} else {
		funcName = runtime.FuncForPC(pc).Name()
	}
	this.record.Pathname = file
	this.record.Line = line
	this.record.Func = funcName

	this.record.Prefix = ""
	this.record.Contexts = this.record.Contexts[:0]
	this.record.Marked = false

	return &this.record
}
