package writer

import (
	"github.com/gxlog/gxlog"
)

var gNullWriter = WriterFunc(func([]byte, *gxlog.Record) {})

func NullWriter() gxlog.Writer {
	return gNullWriter
}
