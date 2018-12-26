package writer

import "github.com/gxlog/gxlog"

// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(bs []byte, record *gxlog.Record)

func (self Func) Write(bs []byte, record *gxlog.Record) {
	self(bs, record)
}
