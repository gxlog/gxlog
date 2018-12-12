package writer

import "github.com/gxlog/gxlog"

type WriterFunc func(bs []byte, record *gxlog.Record)

func (self WriterFunc) Write(bs []byte, record *gxlog.Record) {
	self(bs, record)
}
