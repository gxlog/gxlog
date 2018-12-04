package gxlog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	cCallDepth        = 4
	cBatchRecordCount = 16
)

type logger struct {
	level       LogLevel
	exitOnFatal bool

	linkSlots    [MaxLinkSlot]*link
	compactSlots []*link
	records      []Record

	lock sync.Mutex
}

func (this *logger) Log(calldepth int, level LogLevel, actions []Action, args []interface{}) {
	thisLevel, exitOnFatal := this.getLevelAndExitOnFatal()
	if thisLevel <= level {
		this.write(calldepth, level, actions, fmt.Sprint(args...))
	}
	if exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Logf(calldepth int, level LogLevel, actions []Action,
	fmtstr string, args []interface{}) {
	thisLevel, exitOnFatal := this.getLevelAndExitOnFatal()
	if thisLevel <= level {
		this.write(calldepth, level, actions, fmt.Sprintf(fmtstr, args...))
	}
	if exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Panic(actions []Action, args []interface{}) {
	msg := fmt.Sprint(args...)
	if this.GetLevel() <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(actions []Action, fmtstr string, args []interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	if this.GetLevel() <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) Time(actions []Action, args []interface{}) func() {
	done := func() {}
	if this.GetLevel() <= LevelTrace {
		done = this.genDone(actions, fmt.Sprint(args...))
	}
	return done
}

func (this *logger) Timef(actions []Action, fmtstr string, args []interface{}) func() {
	done := func() {}
	if this.GetLevel() <= LevelTrace {
		done = this.genDone(actions, fmt.Sprintf(fmtstr, args...))
	}
	return done
}

func (this *logger) GetLevel() (level LogLevel) {
	this.lock.Lock()
	level = this.level
	this.lock.Unlock()

	return level
}

func (this *logger) SetLevel(level LogLevel) {
	this.lock.Lock()
	this.level = level
	this.lock.Unlock()
}

func (this *logger) GetExitOnFatal() (ok bool) {
	this.lock.Lock()
	ok = this.exitOnFatal
	this.lock.Unlock()

	return ok
}

func (this *logger) SetExitOnFatal(ok bool) {
	this.lock.Lock()
	this.exitOnFatal = ok
	this.lock.Unlock()
}

func (this *logger) getLevelAndExitOnFatal() (level LogLevel, exitOnFatal bool) {
	this.lock.Lock()
	level = this.level
	exitOnFatal = this.exitOnFatal
	this.lock.Unlock()

	return level, exitOnFatal
}

func (this *logger) write(calldepth int, level LogLevel, actions []Action, msg string) {
	file, line, funcName := getRuntime(calldepth + cCallDepth)

	this.lock.Lock()

	record := this.getRecord()
	record.Time = time.Now()
	record.Level = level
	record.Pathname = file
	record.Line = line
	record.Func = funcName
	record.Msg = msg

	for _, action := range actions {
		action(record)
	}

	for _, lnk := range this.compactSlots {
		lnk.writer.Write(lnk.formatter.Format(record), record)
	}

	this.lock.Unlock()
}

func (this *logger) genDone(actions []Action, msg string) func() {
	now := time.Now()
	return func() {
		costs := time.Since(now)
		if this.GetLevel() <= LevelTrace {
			this.write(-1, LevelTrace, actions, fmt.Sprintf("%s (costs: %v)", msg, costs))
		}
	}
}

func (this *logger) getRecord() *Record {
	if len(this.records) == 0 {
		this.records = make([]Record, cBatchRecordCount)
	}
	record := &this.records[0]
	this.records = this.records[1:]
	return record
}

func getRuntime(calldepth int) (file string, line int, funcName string) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(calldepth)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	} else {
		file = "?file?"
		line = -1
		funcName = "?func?"
	}
	return file, line, funcName
}
