package main

import (
	"fmt"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
	"github.com/gxlog/gxlog/formatter"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/writer"
)

var log = defaults.Logger()

func main() {
	// Only supported on systems that ANSI escape sequences are supported.
	defaults.Formatter().EnableColor()

	log.Info("this will print once")

	// copy Slot0 with the default formatter and wrapper of os.Stderr to Slot1
	log.CopySlot(gxlog.Slot1, gxlog.Slot0)
	log.Info("this will print twice")

	log.SetSlotFormatter(gxlog.Slot1, json.New(json.NewConfig()))
	log.Info("this will print in text format and json format")

	log.SwapSlot(gxlog.Slot0, gxlog.Slot1)
	log.Info("json first and then text")

	log.Unlink(gxlog.Slot0)
	log.Infof("busy slots: %v, free slots: %v", log.BusySlots(), log.FreeSlots())

	log.SetSlotLevel(gxlog.Slot1, gxlog.LevelWarn)
	log.Info("this will not print")
	log.Warn("this will print")

	log.SetSlotLevel(gxlog.Slot1, gxlog.LevelTrace)
	// ATTENTION: DO NOT call methods of logger in formatter, writer or filter
	//   in the current goroutine, or it will deadlock.
	hook := formatter.Func(func(record *gxlog.Record) []byte {
		// log.Info("deadlock")
		fmt.Println("hooks:", record.Msg)
		return nil
	})
	filter := func(record *gxlog.Record) bool {
		return record.Aux.Marked
	}
	// link at Slot0 will overwrite the current link at Slot0 if any
	// If the log level is not lower than WARN and the log is marked, the hook
	//   will be called.
	// use writer.Null() instead of nil, or it will panic
	log.Link(gxlog.Slot0, hook, writer.Null(), gxlog.LevelWarn, filter)
	log.WithMark(true).Info("marked, but info")
	log.Error("error, but not marked")
	log.WithMark(true).Warn("warn and marked")
}
