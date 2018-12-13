package gxlog

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	cCallDepth = 4
)

type Filter func(*Record) bool

type logger struct {
	level      Level
	trackLevel Level
	exitLevel  Level
	filter     Filter
	limit      bool

	slots [MaxSlot]*link

	countMap map[locator]int64

	lock sync.Mutex
}

func (this *logger) Log(calldepth int, level Level, attr *attribute, args []interface{}) {
	logLevel, trackLevel, exitLevel := this.levels()
	if logLevel <= level {
		if trackLevel <= level {
			args = append(args, "\n", string(debug.Stack()))
		}
		this.write(calldepth, level, attr, fmt.Sprint(args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

func (this *logger) Logf(calldepth int, level Level, attr *attribute,
	fmtstr string, args []interface{}) {
	logLevel, trackLevel, exitLevel := this.levels()
	if logLevel <= level {
		if trackLevel <= level {
			fmtstr += "\n%s"
			args = append(args, debug.Stack())
		}
		this.write(calldepth, level, attr, fmt.Sprintf(fmtstr, args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

func (this *logger) Panic(attr *attribute, args []interface{}) {
	msg := fmt.Sprint(args...)
	if this.Level() <= LevelFatal {
		this.write(0, LevelFatal, attr, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(attr *attribute, fmtstr string, args []interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	if this.Level() <= LevelFatal {
		this.write(0, LevelFatal, attr, msg)
	}
	panic(msg)
}

func (this *logger) Time(attr *attribute, args []interface{}) func() {
	if this.Level() <= LevelTrace {
		return this.genDone(attr, fmt.Sprint(args...))
	}
	return func() {}
}

func (this *logger) Timef(attr *attribute, fmtstr string, args []interface{}) func() {
	if this.Level() <= LevelTrace {
		return this.genDone(attr, fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (this *logger) Level() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.level
}

func (this *logger) SetLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.level = level
}

func (this *logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.filter
}

func (this *logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.filter = filter
}

func (this *logger) TrackLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.trackLevel
}

func (this *logger) SetTrackLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.trackLevel = level
}

func (this *logger) ExitLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.exitLevel
}

func (this *logger) SetExitLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.exitLevel = level
}

func (this *logger) Limit() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.limit
}

func (this *logger) SetLimit(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.limit = ok
}

func (this *logger) levels() (Level, Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.level, this.trackLevel, this.exitLevel
}

func (this *logger) write(calldepth int, level Level, attr *attribute, msg string) {
	file, line, pkg, fn := getPosInfo(calldepth + cCallDepth)

	this.lock.Lock()
	defer this.lock.Unlock()

	record := &Record{
		Time:  time.Now(),
		Level: level,
		File:  file,
		Line:  line,
		Pkg:   pkg,
		Func:  fn,
		Msg:   msg,
		Aux:   attr.aux,
	}
	if this.filter != nil && !this.filter(record) {
		return
	}
	if this.limit {
		if attr.countLimiter != nil && !attr.countLimiter(record) {
			return
		}
	}
	for _, link := range this.slots {
		if link != nil && link.level <= level {
			if link.filter == nil || link.filter(record) {
				link.writer.Write(link.formatter.Format(record), record)
			}
		}
	}
}

func (this *logger) genDone(attr *attribute, msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		if this.Level() <= LevelTrace {
			this.write(-1, LevelTrace, attr, fmt.Sprintf("%s (cost: %v)", msg, cost))
		}
	}
}

func getPosInfo(calldepth int) (file string, line int, pkg, fn string) {
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
