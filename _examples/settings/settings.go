package main

import (
	"strings"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/logger"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
	// gxlog.Formatter returns the default Formatter in Slot0.
	// Coloring is only supported on systems that ANSI escape sequences
	// are supported.
	gxlog.Formatter().EnableColor()

	testConfig()
	testFilterLogic()
}

func testConfig() {
	log.Infof("config: %#v", log.Config())

	log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).
		Fatal("fatal before updating the config")
	log.UpdateConfig(func(config logger.Config) logger.Config {
		// Do NOT call any method of the Logger in the function,
		// or it may deadlock.
		// Disable Prefix, StaticContext and Mark, then their value will always
		// be the zero value of their type.
		config.Disabled |= (logger.Prefix | logger.StaticContext | logger.Mark)
		// Disable auto backtracking
		config.TrackLevel = iface.Off
		return config
	})
	log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).
		Fatal("fatal after updating the config")
}

func testFilterLogic() {
	log.SetFilter(logger.Or(important, logger.And(useful, interesting)))
	log.Error("error") // this will be output
	log.Warn("warn")
	log.Trace("trace, funny")
	log.Info("info, funny") // this will be output
}

func important(record *iface.Record) bool {
	return record.Level >= iface.Error
}

func useful(record *iface.Record) bool {
	return record.Level >= iface.Info
}

func interesting(record *iface.Record) bool {
	return strings.Contains(record.Msg, "funny")
}
