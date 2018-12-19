package formatter

import (
	"github.com/gxlog/gxlog"
)

var gNullFormatter = Func(func(*gxlog.Record) []byte { return nil })

func Null() gxlog.Formatter {
	return gNullFormatter
}
