package null

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter"
)

var gNullFormatter = formatter.FormatterFunc(func(*gxlog.Record) []byte { return nil })

func Formatter() gxlog.Formatter {
	return gNullFormatter
}
