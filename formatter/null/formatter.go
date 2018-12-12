package null

import (
	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter"
)

var gNullFormatter = formatter.FormatterFunc(func(*gxlog.Record) []byte { return nil })

func Formatter() gxlog.Formatter {
	return gNullFormatter
}
