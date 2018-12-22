package gxlog

import (
	"fmt"
	"os"
	"path/filepath"
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
	config Config

	slots [MaxSlot]*link

	lock sync.Mutex
}

func (this *logger) Log(calldepth int, level Level, attr *attribute, args ...interface{}) {
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
	fmtstr string, args ...interface{}) {
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

func (this *logger) Panic(attr *attribute, args ...interface{}) {
	msg := fmt.Sprint(args...)
	logLevel, panicLevel := this.panicLevel()
	if logLevel <= panicLevel {
		this.write(0, panicLevel, attr, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(attr *attribute, fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	logLevel, panicLevel := this.panicLevel()
	if logLevel <= panicLevel {
		this.write(0, panicLevel, attr, msg)
	}
	panic(msg)
}

func (this *logger) Time(attr *attribute, args ...interface{}) func() {
	logLevel, timeLevel := this.timeLevel()
	if logLevel <= timeLevel {
		return this.genDone(attr, fmt.Sprint(args...))
	}
	return func() {}
}

func (this *logger) Timef(attr *attribute, fmtstr string, args ...interface{}) func() {
	logLevel, timeLevel := this.timeLevel()
	if logLevel <= timeLevel {
		return this.genDone(attr, fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (this *logger) levels() (Level, Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.TrackLevel, this.config.ExitLevel
}

func (this *logger) timeLevel() (Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.TimeLevel
}

func (this *logger) panicLevel() (Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.PanicLevel
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
	}

	if !this.filter(record, attr) {
		return
	}

	this.attachAux(record, attr)

	for _, link := range this.slots {
		if link != nil && link.Level <= level {
			if link.Filter == nil || link.Filter(record) {
				link.Writer.Write(link.Formatter.Format(record), record)
			}
		}
	}
}

func (this *logger) filter(record *Record, attr *attribute) bool {
	if this.config.Filter != nil && !this.config.Filter(record) {
		return false
	}
	if this.config.Limit {
		if attr.CountLimiter != nil && !attr.CountLimiter(record) {
			return false
		}
		if attr.TimeLimiter != nil && !attr.TimeLimiter(record) {
			return false
		}
	}
	return true
}

func (this *logger) attachAux(record *Record, attr *attribute) {
	if this.config.Prefix {
		record.Aux.Prefix = attr.Prefix
	}
	if this.config.Context {
		// the len and cap of attr.Contexts are equal. next appending will reallocate memory
		record.Aux.Contexts = attr.Contexts
		if this.config.Dynamic {
			for _, context := range attr.DynamicContexts {
				record.Aux.Contexts = append(record.Aux.Contexts, Context{
					Key:   fmt.Sprint(context.Key),
					Value: fmt.Sprint(context.Value(context.Key)),
				})
			}
		}
	}
	if this.config.Mark {
		record.Aux.Marked = attr.Marked
	}
}

func (this *logger) genDone(attr *attribute, msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		logLevel, timeLevel := this.timeLevel()
		if logLevel <= timeLevel {
			this.write(-1, timeLevel, attr, fmt.Sprintf("%s (cost: %v)", msg, cost))
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
	return filepath.ToSlash(file), line, pkg, fn
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
