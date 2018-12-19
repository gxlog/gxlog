package writer

import "github.com/gxlog/gxlog"

type Func func(bs []byte, record *gxlog.Record)

func (self Func) Write(bs []byte, record *gxlog.Record) {
	self(bs, record)
}
