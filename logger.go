package gxlog

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	cCallDepth = 4
)

type logger struct {
	level       LogLevel
	exitOnFatal bool

	linkSlots    [MaxLinkSlot]*link
	compactSlots []*link
	gatherer     gatherer

	lock sync.Mutex
}

func (this *logger) Log(calldepth int, level LogLevel, actions []Action, args []interface{}) {
	this.lock.Lock()

	if this.level <= level {
		this.write(calldepth, level, actions, fmt.Sprint(args...))
	}
	if this.exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}

	this.lock.Unlock()
}

func (this *logger) Logf(calldepth int, level LogLevel, actions []Action,
	fmtstr string, args []interface{}) {
	this.lock.Lock()

	if this.level <= level {
		this.write(calldepth, level, actions, fmt.Sprintf(fmtstr, args...))
	}
	if this.exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}

	this.lock.Unlock()
}

func (this *logger) Panic(actions []Action, args []interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	msg := fmt.Sprint(args...)
	if this.level <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(actions []Action, fmtstr string, args []interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	msg := fmt.Sprintf(fmtstr, args...)
	if this.level <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) Time(actions []Action, args []interface{}) func() {
	done := func() {}

	this.lock.Lock()

	if this.level <= LevelTrace {
		done = this.genDone(actions, fmt.Sprint(args...))
	}

	this.lock.Unlock()

	return done
}

func (this *logger) Timef(actions []Action, fmtstr string, args []interface{}) func() {
	done := func() {}

	this.lock.Lock()

	if this.level <= LevelTrace {
		done = this.genDone(actions, fmt.Sprintf(fmtstr, args...))
	}

	this.lock.Unlock()

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

func (this *logger) write(calldepth int, level LogLevel, actions []Action, msg string) {
	record := this.gatherer.Gather(calldepth+cCallDepth, level, msg)
	for _, action := range actions {
		action(record)
	}
	for _, lnk := range this.compactSlots {
		lnk.writer.Write(lnk.formatter.Format(record), record)
	}
}

func (this *logger) genDone(actions []Action, msg string) func() {
	now := time.Now()
	return func() {
		costs := time.Since(now)
		this.lock.Lock()
		this.write(-1, LevelTrace, actions, fmt.Sprintf("%s (costs: %v)", msg, costs))
		this.lock.Unlock()
	}
}
