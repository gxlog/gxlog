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
	config Config

	slots [MaxSlot]*link

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

func (this *logger) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	return &copyConfig
}

func (this *logger) SetConfig(config *Config) {
	if config == nil {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	this.config = *config
}

func (this *logger) UpdateConfig(fn func(*Config)) {
	if fn == nil {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	fn(&copyConfig)
	this.config = copyConfig
}

func (this *logger) Level() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level
}

func (this *logger) SetLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Level = level
}

func (this *logger) TrackLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TrackLevel
}

func (this *logger) SetTrackLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TrackLevel = level
}

func (this *logger) ExitLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.ExitLevel
}

func (this *logger) SetExitLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.ExitLevel = level
}

func (this *logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Filter
}

func (this *logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Filter = filter
}

func (this *logger) Prefix() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Prefix
}

func (this *logger) SetPrefix(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Prefix = ok
}

func (this *logger) Context() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Context
}

func (this *logger) SetContext(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Context = ok
}

func (this *logger) Mark() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Mark
}

func (this *logger) SetMark(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Mark = ok
}

func (this *logger) Limit() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Limit
}

func (this *logger) SetLimit(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Limit = ok
}

func (this *logger) levels() (Level, Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.TrackLevel, this.config.ExitLevel
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
	if this.config.Filter != nil && !this.config.Filter(record) {
		return
	}
	if this.config.Limit {
		if attr.countLimiter != nil && !attr.countLimiter(record) {
			return
		}
		if attr.timeLimiter != nil && !attr.timeLimiter(record) {
			return
		}
	}
	if this.config.Prefix {
		record.Aux.Prefix = attr.prefix
	}
	if this.config.Context {
		// slicing to set capacity to length, force next appending to reallocate memory
		record.Aux.Contexts = attr.contexts[:len(attr.contexts):len(attr.contexts)]
		for _, context := range attr.dynamicContexts {
			record.Aux.Contexts = append(record.Aux.Contexts, Context{
				Key:   fmt.Sprint(context.key),
				Value: fmt.Sprint(context.value(context.key)),
			})
		}
	}
	if this.config.Mark {
		record.Aux.Marked = attr.marked
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
