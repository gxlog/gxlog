package main

import (
	"fmt"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/logger"
	"github.com/gxlog/gxlog/writer"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
	// gxlog.Formatter returns the default Formatter in Slot0.
	// Coloring is only supported on systems that ANSI escape sequences
	// are supported.
	gxlog.Formatter().EnableColoring()

	testSlots()
	testSlotsLevel()
}

func testSlots() {
	log.Info("this will be printed once")

	// Copy the Formatter, Writer, Filter and Level of Slot0 to Slot1.
	log.CopySlot(logger.Slot1, logger.Slot0)
	log.Info("this will be printed twice")

	log.SetSlotFormatter(logger.Slot1, json.New(json.Config{}))
	log.Info("this will be printed in text format and json format")

	log.SwapSlot(logger.Slot0, logger.Slot1)
	log.Info("json first and then text")

	// Copy the Formatter, Writer, Filter and Level of Slot1 to Slot0 and then
	// set the Formatter to formatter.Null(), Writer to writer.Null(), Filter to
	// nil and Level to Off of Slot1.
	log.MoveSlot(logger.Slot0, logger.Slot1)
}

func testSlotsLevel() {
	log.SetSlotLevel(logger.Slot0, iface.Warn)
	log.Info("this will not be printed")
	log.Warn("this will be printed")

	log.SetSlotLevel(logger.Slot0, iface.Trace)
	// A Formatter or Writer can act as a hook.
	// ATTENTION: Do NOT call any method of the Logger in a Formatter, Writer
	// or Filter, or it may deadlock.
	hook := formatter.Func(func(record *iface.Record) []byte {
		// log.Info("deadlock")
		fmt.Println("hooks:", record.Msg)
		return nil
	})
	filter := func(record *iface.Record) bool {
		return record.Aux.Marked
	}
	// Link at Slot0 will overwrite the current link at Slot0.
	// Use formatter.Null() instead of a nil Formatter and writer.Null()
	// instead of a nil Writer, or it will panic.
	// If the Level of a log is NOT lower than Warn and it is marked, the hook
	// will be called.
	log.Link(logger.Slot0, hook, writer.Null(), iface.Warn, filter)
	log.WithMark(true).Info("marked, but info")
	log.Error("error, but not marked")
	log.WithMark(true).Warn("warn and marked")
}
