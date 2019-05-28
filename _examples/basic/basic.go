package main

import (
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/iface"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
	// gxlog.Formatter returns the default Formatter in Slot0.
	// Coloring is only supported on systems that ANSI escape sequences
	// are supported.
	gxlog.Formatter().EnableColoring()

	testLevel()
	// testPanic()
	testTime()
	testLog()
	testLogError()
	testLogErrorf()
}

func testLevel() {
	log.Trace("test Trace")
	log.Tracef("%s", "test Tracef")
	log.Debug("test Debug")
	log.Debugf("%s", "test Debugf")
	log.Info("test Info")
	log.Infof("%s", "test Infof")
	log.Warn("test Warn")
	log.Warnf("%s", "test Warnf")
	log.Error("test Error")
	log.Errorf("%s", "test Errorf")
	// Fatal and Fatalf will output the stack of current goroutine by default.
	log.Fatal("test Fatal")
	log.Fatalf("%s", "test Fatalf")
}

func testPanic() {
	// The default Level of Panic or Panicf is Fatal.
	// It will always panic when Panic or Panicf is called no matter at which
	// level the Logger is.
	log.Panic("test Panic")
	log.Panicf("%s", "test Panicf")
}

func testTime() {
	// Time or Timef returns a function. When the function is called, it outputs
	// the log as well as the time cost since the call of Time or Timef.
	// The default Level of Time or Timef is Trace.
	done := log.Timing("test Time")
	time.Sleep(200 * time.Millisecond)
	done()
	// Time or Timef works well with defer.
	// Notice the last empty pair of parentheses.
	defer log.Timingf("%s", "test Timef")()
	time.Sleep(400 * time.Millisecond)
}

func testLog() {
	// The calldepth can be specified in Log and Logf. That is useful when you
	// are customizing your own log wrapper function.
	log.Log(0, iface.Info, "test Log")
	log.Logf(1, iface.Warn, "%s: %d", "test Logf", 1)
	log.Logf(-1, iface.Warn, "%s: %d", "test Logf", -1)
}

func testLogError() error {
	// LogError outputs a log and call errors.New to generate an error.
	return log.LogError(iface.Error, "an error")
}

func testLogErrorf() error {
	// LogErrorf outputs a log and call fmt.Errorf to generate an error.
	return log.LogErrorf(iface.Error, "%s", "another error")
}
