package null

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/writer"
)

var Writer = writer.WriterFunc(func([]byte, *gxlog.Record) {})
