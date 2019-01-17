package main

import (
	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testCustomFormatter()
	testTextFormatter()
	testJSONFormatter()
}

func testCustomFormatter() {
	fn := formatter.Func(func(record *iface.Record) []byte {
		return []byte(record.Msg + "\n")
	})
	log.SetSlotFormatter(logger.Slot0, fn)
	log.Info("a simple formatter that just returns the msg of a record")
}

func testTextFormatter() {
	// By default, Trace, Debug and Info map to Green, Warn maps to Yellow,
	// Error and Fatal map to Red, marked logs map to Magenta.
	textFmt := text.New(text.Config{
		// Coloring is only supported on systems that ANSI escape sequences
		// are supported.
		EnableColor: true,
		Header:      text.CompactHeader,
	})
	log.SetSlotFormatter(logger.Slot0, textFmt)
	log.Trace("green")
	log.Warn("yellow")
	log.Error("red")
	log.WithMark(true).Error("magenta")

	// update settings
	textFmt.SetHeader(text.FullHeader)
	textFmt.SetColor(iface.Trace, text.Blue)
	textFmt.MapColors(map[iface.Level]text.Color{
		iface.Warn:  text.Red,
		iface.Error: text.Magenta,
	})
	textFmt.SetMarkedColor(text.White)
	log.Trace("blue")
	log.Warn("red")
	log.Error("magenta")
	log.WithMark(true).Error("white")

	header := "{{time:time}} {{level:char}} {{file:2%q}}:{{line:%05d}} {{msg:%20s}}\n"
	textFmt.SetHeader(header)
	textFmt.DisableColor()
	log.Trace("default color")
}

func testJSONFormatter() {
	jsonFmt := json.New(json.Config{
		// Only the last segment of the File field will be formatted.
		FileSegs: 1,
	})
	log.SetSlotFormatter(logger.Slot0, jsonFmt)
	log.Trace("json")

	// update settings
	jsonFmt.UpdateConfig(func(config json.Config) json.Config {
		// Do NOT call any method of the Formatter or the Logger in the function,
		// or it may deadlock.
		config.OmitEmpty = json.Aux
		config.Omit = json.Pkg | json.Func
		return config
	})
	log.Trace("json updated")
	log.WithContext("ah", "ha").Trace("json with context")
}
