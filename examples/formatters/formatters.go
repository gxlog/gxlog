package main

import (
	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
	"github.com/gxlog/gxlog/formatter"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/formatter/text"
)

var log = defaults.Logger()

func main() {
	// custom formatter function
	fn := formatter.Func(func(record *gxlog.Record) []byte {
		return append([]byte(record.Msg), '\n')
	})
	log.SetSlotFormatter(gxlog.Slot0, fn)
	log.Info("a simple formatter that just returns record.Msg")

	// text formatter
	// the default color mapping is Trace, Debug and Info to Green, Warn to
	//   Yellow, Error and Fatal to Red and marked logs to Magenta no matter
	//   at which level they are.
	textFmt := text.New(text.NewConfig().
		// Only supported on systems that ANSI escape sequences are supported.
		WithEnableColor(true).
		WithHeader(text.CompactHeader))
	log.SetSlotFormatter(gxlog.Slot0, textFmt)
	log.Trace("green")
	log.Warn("yellow")
	log.Error("red")
	log.WithMark(true).Error("magenta")

	// update settings of the text formatter
	textFmt.SetHeader(text.DefaultHeader)
	textFmt.SetColor(gxlog.Trace, text.Blue)
	textFmt.MapColors(map[gxlog.Level]text.Color{
		gxlog.Warn:  text.Red,
		gxlog.Error: text.Magenta,
	})
	textFmt.SetMarkedColor(text.White)
	log.Trace("blue")
	log.Warn("red")
	log.Error("magenta")
	log.WithMark(true).Error("white")

	// custom header of text formatter
	textFmt.SetHeader("{{time:time}} {{level:char}} {{file:2%q}}:{{line:%05d}} {{msg:%20s}}\n")
	textFmt.DisableColor()
	log.Trace("default color")

	// json formatter, with the config that only the last segment of the File
	//   field will be formatted
	jsonFmt := json.New(json.NewConfig().WithFileSegs(1))
	log.SetSlotFormatter(gxlog.Slot0, jsonFmt)
	log.Trace("json")

	// update settings of the json formatter
	jsonFmt.UpdateConfig(func(config json.Config) json.Config {
		// Do NOT call methods of the json formatter, or it will deadlock.
		config.OmitEmpty = json.Aux
		config.Omit = json.Pkg | json.Func
		return config
	})
	log.Trace("json updated")
	log.WithContext("ah", "ha").Trace("json with contexts")
}
