package main

import (
	"strings"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func main() {
	// Only supported on systems that ANSI escape sequences are supported.
	defaults.Formatter().EnableColor()

	log.Infof("config: %#v", log.Config())

	log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).Fatal("fatal before update")
	log.UpdateConfig(func(config gxlog.Config) gxlog.Config {
		// disable prefix, contexts and mark
		// these attributes of records will always be the zero value of their type
		config.Flags &^= (gxlog.Prefix | gxlog.Contexts | gxlog.Mark)
		// disable the auto backtracking
		config.TrackLevel = gxlog.Off
		return config
	})
	log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).Fatal("fatal after update")

	// demonstrates the filter logic
	log.SetFilter(gxlog.Or(important, gxlog.And(useful, interesting)))
	log.Error("error") // this will print
	log.Warn("warn")
	log.Trace("trace, funny")
	log.Info("info, funny") // this will print
}

func important(record *gxlog.Record) bool {
	return record.Level >= gxlog.Error
}

func useful(record *gxlog.Record) bool {
	return record.Level >= gxlog.Info
}

func interesting(record *gxlog.Record) bool {
	return strings.Contains(record.Msg, "funny")
}
