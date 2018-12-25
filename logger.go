package gxlog

import (
	"errors"
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
	cCallDepth = 3
)

type Logger struct {
	config *Config
	slots  []link

	countMap map[locator]int64
	timeMap  map[locator]*timeQueue

	attr copyOnWrite

	lock *sync.Mutex
}

func New(config *Config) *Logger {
	copyConfig := *config
	logger := &Logger{
		config:   &copyConfig,
		countMap: make(map[locator]int64, cMapInitCap),
		timeMap:  make(map[locator]*timeQueue, cMapInitCap),
		lock:     new(sync.Mutex),
	}
	logger.initSlots()
	return logger
}

func (this *Logger) Trace(args ...interface{}) {
	this.Log(1, Trace, args...)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.Logf(1, Trace, fmtstr, args...)
}

func (this *Logger) Debug(args ...interface{}) {
	this.Log(1, Debug, args...)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.Logf(1, Debug, fmtstr, args...)
}

func (this *Logger) Info(args ...interface{}) {
	this.Log(1, Info, args...)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.Logf(1, Info, fmtstr, args...)
}

func (this *Logger) Warn(args ...interface{}) {
	this.Log(1, Warn, args...)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.Logf(1, Warn, fmtstr, args...)
}

func (this *Logger) Error(args ...interface{}) {
	this.Log(1, Error, args...)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.Logf(1, Error, fmtstr, args...)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.Log(1, Fatal, args...)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.Logf(1, Fatal, fmtstr, args...)
}

func (this *Logger) LogError(level Level, text string) error {
	this.Log(1, level, text)
	return errors.New(text)
}

func (this *Logger) LogErrorf(level Level, fmtstr string, args ...interface{}) error {
	err := fmt.Errorf(fmtstr, args...)
	this.Log(1, level, err.Error())
	return err
}

func (this *Logger) Log(calldepth int, level Level, args ...interface{}) {
	logLevel, trackLevel, exitLevel := this.levels()
	if logLevel <= level {
		if trackLevel <= level {
			args = append(args, "\n", string(debug.Stack()))
		}
		this.write(calldepth, level, fmt.Sprint(args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

func (this *Logger) Logf(calldepth int, level Level, fmtstr string, args ...interface{}) {
	logLevel, trackLevel, exitLevel := this.levels()
	if logLevel <= level {
		if trackLevel <= level {
			fmtstr += "\n%s"
			args = append(args, debug.Stack())
		}
		this.write(calldepth, level, fmt.Sprintf(fmtstr, args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

func (this *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logLevel, panicLevel := this.panicLevel()
	if logLevel <= panicLevel {
		this.write(0, panicLevel, msg)
	}
	panic(msg)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	logLevel, panicLevel := this.panicLevel()
	if logLevel <= panicLevel {
		this.write(0, panicLevel, msg)
	}
	panic(msg)
}

func (this *Logger) Time(args ...interface{}) func() {
	logLevel, timeLevel := this.timeLevel()
	if logLevel <= timeLevel {
		return this.genDone(fmt.Sprint(args...))
	}
	return func() {}
}

func (this *Logger) Timef(fmtstr string, args ...interface{}) func() {
	logLevel, timeLevel := this.timeLevel()
	if logLevel <= timeLevel {
		return this.genDone(fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (this *Logger) levels() (Level, Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.TrackLevel, this.config.ExitLevel
}

func (this *Logger) timeLevel() (Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.TimeLevel
}

func (this *Logger) panicLevel() (Level, Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level, this.config.PanicLevel
}

func (this *Logger) write(calldepth int, level Level, msg string) {
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

	if !this.filter(record) {
		return
	}

	this.attachAux(record)

	for _, link := range this.slots {
		if link.Level <= level {
			if link.Filter == nil || link.Filter(record) {
				var bs []byte
				if link.Formatter != nil {
					bs = link.Formatter.Format(record)
				}
				if link.Writer != nil {
					link.Writer.Write(bs, record)
				}
			}
		}
	}
}

func (this *Logger) filter(record *Record) bool {
	if this.config.Filter != nil && !this.config.Filter(record) {
		return false
	}
	if this.config.Flags&Limit != 0 {
		if this.attr.CountLimiter != nil && !this.attr.CountLimiter(record) {
			return false
		}
		if this.attr.TimeLimiter != nil && !this.attr.TimeLimiter(record) {
			return false
		}
	}
	return true
}

func (this *Logger) attachAux(record *Record) {
	if this.config.Flags&Prefix != 0 {
		record.Aux.Prefix = this.attr.Prefix
	}
	if this.config.Flags&Contexts != 0 {
		// The len and cap of this.attr.Contexts are equal,
		//   next appending will reallocate memory
		record.Aux.Contexts = this.attr.Contexts
		if this.config.Flags&DynamicContexts != 0 {
			for _, context := range this.attr.DynamicContexts {
				record.Aux.Contexts = append(record.Aux.Contexts, Context{
					Key:   fmt.Sprint(context.Key),
					Value: fmt.Sprint(context.Value(context.Key)),
				})
			}
		}
	}
	if this.config.Flags&Mark != 0 {
		record.Aux.Marked = this.attr.Marked
	}
}

func (this *Logger) genDone(msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		logLevel, timeLevel := this.timeLevel()
		if logLevel <= timeLevel {
			this.write(0, timeLevel, fmt.Sprintf("%s (cost: %v)", msg, cost))
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
