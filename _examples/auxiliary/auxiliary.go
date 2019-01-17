package main

import (
	"math/rand"
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/logger"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// gxlog.Formatter returns the default Formatter in Slot0.
	// Coloring is only supported on systems that ANSI escape sequences
	// are supported.
	gxlog.Formatter().EnableColor()

	testAuxiliary()
	testDynamicContext()
	testLimitation()
}

func testAuxiliary() {
	// Logs with mark will be colorized with Magenta by default.
	// The prefix or mark allow you to highlight some logs temporarily
	// while you are debugging.
	log.WithPrefix("**** ").WithMark(true).WithContext("k1", "v1", "k2", "v2").
		Info("prefix, mark and contexts")
	// The original log instance is not altered.
	log.Info("no prefix, mark or contexts")

	// This demonstrates the lexical scope of a log instance:
	func() {
		log := log.WithContext("k3", "v3")
		log.Info("outer enter")
		func() {
			log := log.WithContext("k4", "v4")
			log.Info("inner")
		}()
		log.Info("outer leave")
	}()
}

func testDynamicContext() {
	// All the key-value pairs of dynamic contexts will be appended to the end
	// of static contexts.
	// Dynamic contexts are very useful when you want to print the current value
	// of some variables all the time.
	// ATTENTION: You SHOULD be very careful to concurrency safety or deadlocks
	// with dynamic contexts.
	n := 0
	fn := logger.Dynamic(func(interface{}) interface{} {
		// Do NOT call any method of the Logger in the function,
		// or it may deadlock.
		n++
		return n
	})
	clog := log.WithContext("static", n, "dynamic", fn)
	clog.Info("dynamic one")
	clog.Info("dynamic two")
}

func testLimitation() {
	// THINK TWICE before you decide to limit the output of logs by count or
	// by time, you may miss logs which you need.
	// Only 2 logs will be output per 3 logs.
	for i := 1; i <= 6; i++ {
		log.WithCountLimit(3, 2).Infof("count limited: %d", i)
	}
	// the more efficient way
	llog := log.WithCountLimit(3, 2)
	for i := 7; i <= 12; i++ {
		llog.Infof("efficient count limited: %d", i)
	}
	// NOTICE: The space complexity is O(n), while n is the 2nd argument of
	// WithTimeLimit. Try to specify reasonable duration and limit.
	// At most 3 logs will be output during any interval of 1 second.
	for i := 1; i <= 10; i++ {
		log.WithTimeLimit(time.Second, 3).Infof("time limited: %d", i)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
