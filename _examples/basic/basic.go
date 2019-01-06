package main

import (
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func main() {
	// Only supported on systems that ANSI escape sequences are supported.
	defaults.Formatter().EnableColor()

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

	// The default level of Panic or Panicf is fatal.
	// It will always panic no matter at which level the logger is.
	// log.Panic("test Panic")
	// log.Panicf("%s", "test Panicf")

	// Time and Timef will return a function. When the function is called,
	//   it will output the log as well as the time cost since the call of
	//   Time or Timef.
	// The default level of Time and Timef is trace.
	done := log.Time("test Time")
	time.Sleep(200 * time.Millisecond)
	done()
	// Notice the last empty pair of parentheses.
	defer log.Timef("%s", "test Timef")()
	time.Sleep(400 * time.Millisecond)

	// The calldepth can be specified in Log and Logf. That is useful when
	//   you want to customize your own log helper functions.
	log.Log(0, gxlog.Info, "test Log")
	log.Logf(1, gxlog.Warn, "%s: %d", "test Logf", 1)
	log.Logf(-1, gxlog.Warn, "%s: %d", "test Logf", -1)

	test1()
	test2()
}

func test1() error {
	// LogError will output log and call errors.New to generate an error
	return log.LogError(gxlog.Error, "an error")
}

func test2() error {
	// LogErrorf will output log and call fmt.Errorf to generate an error
	return log.LogErrorf(gxlog.Error, "%s", "another error")
}
