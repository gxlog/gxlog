package formatter

import (
	"github.com/gxlog/gxlog"
)

var gNullFormatter = FormatterFunc(func(*gxlog.Record) []byte { return nil })

func NullFormatter() gxlog.Formatter {
	return gNullFormatter
}
