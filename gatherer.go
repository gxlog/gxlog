package gxlog

import (
	"runtime"
	"time"
)

const (
	cBatchRecordCount = 16
)

type gatherer struct {
	records []Record
}

func (this *gatherer) Gather(calldepth int, level LogLevel, msg string) *Record {
	record := this.getRecord()
	record.Time = time.Now()
	record.Level = level
	record.Msg = msg

	var funcName string
	pc, file, line, ok := runtime.Caller(calldepth)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	} else {
		file = "?file?"
		line = -1
		funcName = "?func?"
	}
	record.Pathname = file
	record.Line = line
	record.Func = funcName

	return record
}

func (this *gatherer) getRecord() *Record {
	if len(this.records) == 0 {
		this.records = make([]Record, cBatchRecordCount)
	}
	record := &this.records[0]
	this.records = this.records[1:]
	return record
}
