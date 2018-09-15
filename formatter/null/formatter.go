package null

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter"
)

var nullFormatter = formatter.FormatterFunc(func(*gxlog.Record) []byte { return nil })

func Formatter() gxlog.Formatter {
	return nullFormatter
}
