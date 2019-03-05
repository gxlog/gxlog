// Package logger implements a concise, functional, flexible and extensible
// logger for Go. It provides many advanced features and has good customizability
// and extensibility. It also aims at easy-to-use.
package logger

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

	"github.com/gxlog/gxlog/iface"
)

const callDepthOffset = 3

// A Logger is a logging framework that contains EIGHT slots. Each Slot contains
// a Formatter and a Writer. A Logger has its own level and filter while each
// Slot has its independent level and filter. Logger calls the Formatter and
// Writer of each Slot in the order from Slot0 to Slot7 when a log is emitted.
//
// All methods of A Logger are concurrency safe.
// A Logger MUST be created with New.
type Logger struct {
	config *Config
	slots  []slotLink
	// store indexes of equivalent formatters, used to avoid redundant formatting
	equivalents [][]int
	countMap    map[locator]int64
	timeMap     map[locator]*timeQueue
	attr        copyOnWrite
	lock        *sync.Mutex
}

// New creates a new Logger with the config.
func New(config Config) *Logger {
	config.setDefaults()
	logger := &Logger{
		config:      &config,
		equivalents: make([][]int, MaxSlot),
		countMap:    make(map[locator]int64, mapInitCap),
		timeMap:     make(map[locator]*timeQueue, mapInitCap),
		lock:        new(sync.Mutex),
	}
	logger.initSlots()
	return logger
}

// Trace calls Log with level Trace to emit a log.
func (log *Logger) Trace(args ...interface{}) {
	log.Log(1, iface.Trace, args...)
}

// Tracef calls Logf with level Trace to emit a log.
func (log *Logger) Tracef(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Trace, fmtstr, args...)
}

// Debug calls Log with level Debug to emit a log.
func (log *Logger) Debug(args ...interface{}) {
	log.Log(1, iface.Debug, args...)
}

// Debugf calls Logf with level Debug to emit a log.
func (log *Logger) Debugf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Debug, fmtstr, args...)
}

// Info calls Log with level Info to emit a log.
func (log *Logger) Info(args ...interface{}) {
	log.Log(1, iface.Info, args...)
}

// Infof calls Logf with level Info to emit a log.
func (log *Logger) Infof(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Info, fmtstr, args...)
}

// Warn calls Log with level Warn to emit a log.
func (log *Logger) Warn(args ...interface{}) {
	log.Log(1, iface.Warn, args...)
}

// Warnf calls Logf with level Warn to emit a log.
func (log *Logger) Warnf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Warn, fmtstr, args...)
}

// Error calls Log with level Error to emit a log.
func (log *Logger) Error(args ...interface{}) {
	log.Log(1, iface.Error, args...)
}

// Errorf calls Logf with level Error to emit a log.
func (log *Logger) Errorf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Error, fmtstr, args...)
}

// Fatal calls Log with level Fatal to emit a log.
func (log *Logger) Fatal(args ...interface{}) {
	log.Log(1, iface.Fatal, args...)
}

// Fatalf calls Logf with level Fatal to emit a log.
func (log *Logger) Fatalf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Fatal, fmtstr, args...)
}

// LogError calls Log to emit a log and calls errors.New to return an error.
// The level MUST be between Trace and Fatal inclusive.
func (log *Logger) LogError(level iface.Level, text string) error {
	log.Log(1, level, text)
	return errors.New(text)
}

// LogErrorf calls Logf to emit a log and calls fmt.Errorf to return an error.
// The level MUST be between Trace and Fatal inclusive.
func (log *Logger) LogErrorf(level iface.Level, fmtstr string,
	args ...interface{}) error {
	err := fmt.Errorf(fmtstr, args...)
	log.Log(1, level, err.Error())
	return err
}

// Log calls the Formatter and Writer in each Slot to format and write a log.
//
// The level MUST be between Trace and Fatal inclusive. If the level is lower
// than the level of Logger, the log will NOT be emitted. If the level is lower
// than the level of a Slot, the Formatter and Writer of the Slot will NOT be
// called. If the level is NOT lower than the track level of Logger, the stack of
// the current goroutine will be output. If the level is NOT lower than the exit
// level of Logger, the Logger will call os.Exit at last.
//
// The callDepth is used to set the offset of stack. It makes sense when you are
// customizing your own log wrapper function. Otherwise, 0 is just ok.
//
// The args are handled in the manner of fmt.Sprint.
//
// ATTENTION: the log may NOT be output when a Writer is in asynchronous mode and
// os.Exit has been called.
func (log *Logger) Log(callDepth int, level iface.Level, args ...interface{}) {
	logLevel, trackLevel, exitLevel := log.levels()
	if logLevel <= level {
		if trackLevel <= level {
			args = append(args, "\n", string(debug.Stack()))
		}
		log.write(callDepth, level, fmt.Sprint(args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

// Logf does the same with Log except that it calls fmt.Sprintf to format a log.
//
// ATTENTION: the log may NOT be output when a Writer is in asynchronous mode and
// os.Exit has been called.
func (log *Logger) Logf(callDepth int, level iface.Level, fmtstr string, args ...interface{}) {
	logLevel, trackLevel, exitLevel := log.levels()
	if logLevel <= level {
		if trackLevel <= level {
			fmtstr += "\n%s"
			args = append(args, debug.Stack())
		}
		log.write(callDepth, level, fmt.Sprintf(fmtstr, args...))
		if exitLevel <= level {
			os.Exit(1)
		}
	}
}

// Panic calls the Formatters and Writers in each Slot to format and write a log
// and panics at last.
//
// The level of the emitted log is the panic level of Logger. If the level is
// lower than the level of Logger, the log will NOT be output. If the level is
// lower than the level of a Slot, the Formatter and Writer of the Slot will NOT
// be called.
//
// The args are handled in the manner of fmt.Sprint.
//
// ATTENTION: the log may NOT be output when a Writer is in asynchronous mode and
// there is NO recovery after panicking.
// Panic never outputs the stack with the log or calls os.Exit.
func (log *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logLevel, panicLevel := log.panicLevel()
	if logLevel <= panicLevel {
		log.write(0, panicLevel, msg)
	}
	panic(msg)
}

// Panicf does the same with Panic except it calls fmt.Sprintf to format a log.
//
// ATTENTION: the log may NOT be output when a Writer is in asynchronous mode and
// there is NO recovery after panicking.
// Panicf never outputs the stack with the log or calls os.Exit.
func (log *Logger) Panicf(fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	logLevel, panicLevel := log.panicLevel()
	if logLevel <= panicLevel {
		log.write(0, panicLevel, msg)
	}
	panic(msg)
}

// Timing returns a function. When the function is called, it outputs the log with
// the time elapsed since the call of Timing.
//
// The level of the emitted log is the timing level of Logger. If the level is
// lower than the level of Logger, the log will NOT be output. If the level is
// lower than the level of a Slot, the Formatter and Writer of the Slot will NOT
// be called.
//
// The args are handled in the manner of fmt.Sprint.
//
// It works well with `defer' and do NOT forget the last empty pair of parentheses.
func (log *Logger) Timing(args ...interface{}) func() {
	logLevel, timingLevel := log.timingLevel()
	if logLevel <= timingLevel {
		return log.genDone(fmt.Sprint(args...))
	}
	return func() {}
}

// Timingf does the same with Timing except it calls fmt.Sprintf to format a log.
//
// It works well with `defer' and do NOT forget the last empty pair of parentheses.
func (log *Logger) Timingf(fmtstr string, args ...interface{}) func() {
	logLevel, timingLevel := log.timingLevel()
	if logLevel <= timingLevel {
		return log.genDone(fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (log *Logger) levels() (iface.Level, iface.Level, iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.TrackLevel, log.config.ExitLevel
}

func (log *Logger) timingLevel() (iface.Level, iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.TimingLevel
}

func (log *Logger) panicLevel() (iface.Level, iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level, log.config.PanicLevel
}

func (log *Logger) write(callDepth int, level iface.Level, msg string) {
	if level < iface.Trace || level > iface.Fatal {
		panic("logger: invalid level")
	}

	file, line, pkg, fn := "", 0, "", ""
	if log.config.Disabled&Runtime == 0 {
		file, line, pkg, fn = getPosInfo(callDepth + callDepthOffset)
	}

	log.lock.Lock()
	defer log.lock.Unlock()

	record := &iface.Record{
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

	var formats [MaxSlot][]byte
	for slot := 0; slot < MaxSlot; slot++ {
		link := &log.slots[slot]
		if link.Level > level {
			continue
		}
		if link.Filter != nil && !link.Filter(record) {
			continue
		}
		format := formats[slot]
		if format == nil {
			format = link.Formatter.Format(record)
			for _, id := range log.equivalents[slot] {
				formats[id] = format
			}
		}
		link.Writer.Write(format, record)
	}
}

func (log *Logger) filter(record *iface.Record) bool {
	if log.config.Filter != nil && !log.config.Filter(record) {
		return false
	}
	if log.config.Disabled&LimitByCount == 0 {
		if log.attr.CountLimiter != nil && !log.attr.CountLimiter(record) {
			return false
		}
	}
	if log.config.Disabled&LimitByTime == 0 {
		if log.attr.TimeLimiter != nil && !log.attr.TimeLimiter(record) {
			return false
		}
	}
	return true
}

func (log *Logger) attachAux(record *iface.Record) {
	if log.config.Disabled&Prefix == 0 {
		record.Aux.Prefix = log.attr.Prefix
	}
	if log.config.Disabled&StaticContext == 0 {
		record.Aux.Contexts = log.attr.Contexts
	}
	if log.config.Disabled&DynamicContext == 0 {
		for _, context := range log.attr.DynamicContexts {
			record.Aux.Contexts = append(record.Aux.Contexts, iface.Context{
				Key:   fmt.Sprint(context.Key),
				Value: fmt.Sprint(context.Value(context.Key)),
			})
		}
	}
	if log.config.Disabled&Mark == 0 {
		record.Aux.Marked = log.attr.Marked
	}
}

func (log *Logger) genDone(msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		logLevel, timingLevel := log.timingLevel()
		if logLevel <= timingLevel {
			log.write(0, timingLevel, fmt.Sprintf("%s (cost: %v)", msg, cost))
		}
	}
}

func getPosInfo(callDepth int) (file string, line int, pkg, fn string) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(callDepth)
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
