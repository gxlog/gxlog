package writer

import (
	"github.com/gxlog/gxlog"
)

var gNullWriter = Func(func([]byte, *gxlog.Record) {})

func Null() gxlog.Writer {
	return gNullWriter
}
