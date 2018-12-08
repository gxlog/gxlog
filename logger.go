package gxlog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

	lock sync.Mutex
}

func (this *logger) Log(calldepth int, level LogLevel, aux *Auxiliary, args []interface{}) {
	thisLevel, exitOnFatal := this.getLevelAndExitOnFatal()
	if thisLevel <= level {
		this.write(calldepth, level, aux, fmt.Sprint(args...))
	}
	if exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Logf(calldepth int, level LogLevel, aux *Auxiliary,
	fmtstr string, args []interface{}) {
	thisLevel, exitOnFatal := this.getLevelAndExitOnFatal()
	if thisLevel <= level {
		this.write(calldepth, level, aux, fmt.Sprintf(fmtstr, args...))
	}
	if exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Panic(aux *Auxiliary, args []interface{}) {
	msg := fmt.Sprint(args...)
	if this.GetLevel() <= LevelFatal {
		this.write(0, LevelFatal, aux, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(aux *Auxiliary, fmtstr string, args []interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	if this.GetLevel() <= LevelFatal {
		this.write(0, LevelFatal, aux, msg)
	}
	panic(msg)
}

func (this *logger) Time(aux *Auxiliary, args []interface{}) func() {
	done := func() {}
	if this.GetLevel() <= LevelTrace {
		done = this.genDone(aux, fmt.Sprint(args...))
	}
	return done
}

func (this *logger) Timef(aux *Auxiliary, fmtstr string, args []interface{}) func() {
	done := func() {}
	if this.GetLevel() <= LevelTrace {
		done = this.genDone(aux, fmt.Sprintf(fmtstr, args...))
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

func (this *logger) write(calldepth int, level LogLevel, aux *Auxiliary, msg string) {
	file, line, pkg, fn := getRuntime(calldepth + cCallDepth)

	this.lock.Lock()

	record := &Record{
		Time:  time.Now(),
		Level: level,
		File:  file,
		Line:  line,
		Pkg:   pkg,
		Func:  fn,
		Msg:   msg,
		Aux:   *aux,
	}
	for _, lnk := range this.compactSlots {
		if lnk.level <= level {
			lnk.writer.Write(lnk.formatter.Format(record), record)
		}
	}

	this.lock.Unlock()
}

func (this *logger) genDone(aux *Auxiliary, msg string) func() {
	now := time.Now()
	return func() {
		costs := time.Since(now)
		if this.GetLevel() <= LevelTrace {
			this.write(-1, LevelTrace, aux, fmt.Sprintf("%s (costs: %v)", msg, costs))
		}
	}
}

func getRuntime(calldepth int) (file string, line int, pkg, fn string) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(calldepth)
	if ok {
		name := runtime.FuncForPC(pc).Name()
		pkg, fn = splitPkgAndFunc(name)
	} else {
		file = "?file?"
		line = -1
		pkg = "?pkg?"
		fn = "?func?"
	}
	return file, line, pkg, fn
}

func splitPkgAndFunc(name string) (string, string) {
	lastSlash := strings.LastIndexByte(name, '/')
	nextDot := strings.IndexByte(name[lastSlash+1:], '.')
	if nextDot < 0 {
		return "?pkg?", "?func?"
	}
	nextDot += (lastSlash + 1)
	return name[:nextDot], name[nextDot+1:]
}
