package main

import (
	"math/rand"
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Only supported on systems that ANSI escape sequences are supported.
	defaults.Formatter().EnableColor()

	// logs with mark will be colorized with magenta by default
	log.WithPrefix("**** ").WithMark(true).WithContext("k1", "v1", "k2", "v2").
		Info("prefix, mark and contexts")
	// the instance of log is left to be unchanged
	log.Info("no prefix, mark or contexts")

	// demonstrates the lexical scope
	func() {
		log := log.WithContext("k3", "v3")
		log.Info("outer enter")
		func() {
			log := log.WithContext("k4", "v4")
			log.Info("inner")
		}()
		log.Info("outer leave")
	}()

	// all the key-value pairs of dynamic contexts will be concatenated to the
	//   end of static contexts
	// dynamic contexts are very useful when you want to print some/all fields
	//   of a struct value all the time.
	// ATTENTION: you should be very careful to concurrency safety or dead
	//   locks with dynamic contexts.
	n := 0
	fn := gxlog.Dynamic(func(interface{}) interface{} {
		n++
		return n
	})
	clog := log.WithContext("static", n, "dynamic", fn)
	clog.Info("dynamic one")
	clog.Info("dynamic two")

	// THINK TWICE before you limit logs output by count or by time,
	//   you may miss logs which you need.
	// only 2 logs will be output per 3 logs
	for i := 1; i <= 6; i++ {
		log.WithCountLimit(3, 2).Infof("count limited: %d", i)
	}
	// the more efficient way
	llog := log.WithCountLimit(3, 2)
	for i := 7; i <= 12; i++ {
		llog.Infof("efficient count limited: %d", i)
	}
	// NOTICE: The space complexity is O(n), while n is the 2nd argument of
	//   WithTimeLimit. Try to specify reasonable duration and limit.
	// at most 3 logs will be output during any interval of 1 second
	for i := 1; i <= 10; i++ {
		log.WithTimeLimit(time.Second, 3).Infof("time limited: %d", i)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
