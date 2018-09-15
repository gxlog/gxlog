package null

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/writer"
)

var nullWriter = writer.WriterFunc(func([]byte, *gxlog.Record) {})

func Writer() gxlog.Writer {
	return nullWriter
}
