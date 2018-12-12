package null

import (
	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/writer"
)

var gNullWriter = writer.WriterFunc(func([]byte, *gxlog.Record) {})

func Writer() gxlog.Writer {
	return gNullWriter
}
