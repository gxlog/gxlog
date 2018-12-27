// Package gxlog implements a concise, functional, flexible and extensible logger
// for Go. It provides many advanced features and has good customizability and
// extensibility. It also aims at easy-to-use.
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

const callDepthOffset = 3

// A Logger is a framework that contains EIGHT slots. Each slot has a formatter
// and a writer. A Logger has its own level and filter, while each slot has its
// independent level and filter. The Logger will call the formatter and the
// writer of each slot in the order from Slot0 to Slot7 if they are not nil.
//
// All methods of A Logger are concurrency safe.
//
// A Logger must be created with New.
type Logger struct {
	config *Config
	slots  []slotLink

	countMap map[locator]int64
	timeMap  map[locator]*timeQueue

	attr copyOnWrite

	lock *sync.Mutex
}

// New creates a new Logger with the config. The config must not be nil.
func New(config *Config) *Logger {
	copyConfig := *config
	logger := &Logger{
		config:   &copyConfig,
		countMap: make(map[locator]int64, mapInitCap),
		timeMap:  make(map[locator]*timeQueue, mapInitCap),
		lock:     new(sync.Mutex),
	}
	logger.initSlots()
	return logger
}

// Trace calls Log with level Trace to output log.
func (log *Logger) Trace(args ...interface{}) {
	log.Log(1, Trace, args...)
}

// Tracef calls Logf with level Trace to output log.
func (log *Logger) Tracef(fmtstr string, args ...interface{}) {
	log.Logf(1, Trace, fmtstr, args...)
}

// Debug calls Log with level Debug to output log.
func (log *Logger) Debug(args ...interface{}) {
	log.Log(1, Debug, args...)
}

// Debugf calls Logf with level Debug to output log.
func (log *Logger) Debugf(fmtstr string, args ...interface{}) {
	log.Logf(1, Debug, fmtstr, args...)
}

// Info calls Log with level Info to output log.
func (log *Logger) Info(args ...interface{}) {
	log.Log(1, Info, args...)
}

// Infof calls Logf with level Info to output log.
func (log *Logger) Infof(fmtstr string, args ...interface{}) {
	log.Logf(1, Info, fmtstr, args...)
}

// Warn calls Log with level Warn to output log.
func (log *Logger) Warn(args ...interface{}) {
	log.Log(1, Warn, args...)
}

// Warnf calls Logf with level Warn to output log.
func (log *Logger) Warnf(fmtstr string, args ...interface{}) {
	log.Logf(1, Warn, fmtstr, args...)
}

// Error calls Log with level Error to output log.
func (log *Logger) Error(args ...interface{}) {
	log.Log(1, Error, args...)
}

// Errorf calls Logf with level Error to output log.
func (log *Logger) Errorf(fmtstr string, args ...interface{}) {
	log.Logf(1, Error, fmtstr, args...)
}

// Fatal calls Log with level Fatal to output log.
func (log *Logger) Fatal(args ...interface{}) {
	log.Log(1, Fatal, args...)
}

// Fatalf calls Logf with level Fatal to output log.
func (log *Logger) Fatalf(fmtstr string, args ...interface{}) {
	log.Logf(1, Fatal, fmtstr, args...)
}

// LogError calls Log to output log and calls errors.New to return an error.
// The level must be between Trace and Fatal inclusive.
func (log *Logger) LogError(level Level, text string) error {
	log.Log(1, level, text)
	return errors.New(text)
}

// LogErrorf calls Logf to output log and calls fmt.Errorf to return an error.
// The level must be between Trace and Fatal inclusive.
func (log *Logger) LogErrorf(level Level, fmtstr string, args ...interface{}) error {
	err := fmt.Errorf(fmtstr, args...)
	log.Log(1, level, err.Error())
	return err
}

// Log calls formatters and writers in slots to output log.
//
// The level must be between Trace and Fatal inclusive. If the level is lower
// than the level of Logger, NO log will output. If the level is lower than
// the level of a slot, the formatter and writer of the slot will NOT be called.
// If the level is NOT lower than the track level of Logger, the stack of the
// current goroutine will be output. If the level is NOT lower than the exit
// level of Logger, it will call os.Exit to exit.
//
// The calldepth is useful when you customize your own log helper functions
// and need to specify the stack level. Otherwise, 0 is just ok.
//
// The args are handled in the manner of fmt.Sprint.
//
// ATTENTION: log may NOT be output in asynchronous mode if os.Exit has been called.
func (log *Logger) Log(calldepth int, level Level, args ...interface{}) {
	logLevel, trackLevel, exitLevel := log.levels()
	if logLevel <= level {
		if trackLevel <= level {
			args = append(args, "\n", string(debug.Stack()))
		}
		log.write(calldepth, level, fmt.Sprint(args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

// Logf does the same with Log except that it calls fmt.Sprintf to do the format.
//
// ATTENTION: log may NOT be output in asynchronous mode if os.Exit has been called.
func (log *Logger) Logf(calldepth int, level Level, fmtstr string, args ...interface{}) {
	logLevel, trackLevel, exitLevel := log.levels()
	if logLevel <= level {
		if trackLevel <= level {
			fmtstr += "\n%s"
			args = append(args, debug.Stack())
		}
		log.write(calldepth, level, fmt.Sprintf(fmtstr, args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

// Panic calls formatters and writers in slots to output log and then panic.
// It will always panic no matter at which level the Logger is.
//
// The level of this call is the panic level of Logger. If the level is lower
// than the level of Logger, NO log will output. If the level is lower than
// the level of a slot, the formatter and writer of the slot will NOT be called.
//
// The args are handled in the manner of fmt.Sprint.
//
// ATTENTION: log may NOT be output in asynchronous mode if there is no recovery.
func (log *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logLevel, panicLevel := log.panicLevel()
	if logLevel <= panicLevel {
		log.write(0, panicLevel, msg)
	}
	panic(msg)
}

// Panicf does the same with Panic except that it calls fmt.Sprintf to do the format.
//
// ATTENTION: log may NOT be output in asynchronous mode if there is no recovery.
func (log *Logger) Panicf(fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	logLevel, panicLevel := log.panicLevel()
	if logLevel <= panicLevel {
		log.write(0, panicLevel, msg)
	}
	panic(msg)
}

// Time will return a function. When the function is called, it will call
// formatters and writers in slots to output the log as well as the time cost
// since the call of Time. It works well with defer, but do not forget the
// last empty pair of parentheses.
//
// The level of this call is the time level of Logger. If the level is lower
// than the level of Logger, NO log will output. If the level is lower than
// the level of a slot, the formatter and writer of the slot will NOT be called.
func (log *Logger) Time(args ...interface{}) func() {
	logLevel, timeLevel := log.timeLevel()
	if logLevel <= timeLevel {
		return log.genDone(fmt.Sprint(args...))
	}
	return func() {}
}

// Timef does the same with Time except that it calls fmt.Sprintf to do the format.
func (log *Logger) Timef(fmtstr string, args ...interface{}) func() {
	logLevel, timeLevel := log.timeLevel()
	if logLevel <= timeLevel {
		return log.genDone(fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (log *Logger) levels() (Level, Level, Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.TrackLevel, log.config.ExitLevel
}

func (log *Logger) timeLevel() (Level, Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.TimeLevel
}

func (log *Logger) panicLevel() (Level, Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.PanicLevel
}

func (log *Logger) write(calldepth int, level Level, msg string) {
	file, line, pkg, fn := getPosInfo(calldepth + callDepthOffset)

	log.lock.Lock()
	defer log.lock.Unlock()

	record := &Record{
		Time:  time.Now(),
		Level: level,
		File:  file,
		Line:  line,
		Pkg:   pkg,
		Func:  fn,
		Msg:   msg,
	}

	if !log.filter(record) {
		return
	}

	log.attachAux(record)

	for _, link := range log.slots {
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

func (log *Logger) filter(record *Record) bool {
	if log.config.Filter != nil && !log.config.Filter(record) {
		return false
	}
	if log.config.Flags&Limit != 0 {
		if log.attr.CountLimiter != nil && !log.attr.CountLimiter(record) {
			return false
		}
		if log.attr.TimeLimiter != nil && !log.attr.TimeLimiter(record) {
			return false
		}
	}
	return true
}

func (log *Logger) attachAux(record *Record) {
	if log.config.Flags&Prefix != 0 {
		record.Aux.Prefix = log.attr.Prefix
	}
	if log.config.Flags&Contexts != 0 {
		// The len and cap of log.attr.Contexts are equal,
		//   next appending will reallocate memory
		record.Aux.Contexts = log.attr.Contexts
		if log.config.Flags&DynamicContexts != 0 {
			for _, context := range log.attr.DynamicContexts {
				record.Aux.Contexts = append(record.Aux.Contexts, Context{
					Key:   fmt.Sprint(context.Key),
					Value: fmt.Sprint(context.Value(context.Key)),
				})
			}
		}
	}
	if log.config.Flags&Mark != 0 {
		record.Aux.Marked = log.attr.Marked
	}
}

func (log *Logger) genDone(msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		logLevel, timeLevel := log.timeLevel()
		if logLevel <= timeLevel {
			log.write(0, timeLevel, fmt.Sprintf("%s (cost: %v)", msg, cost))
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
