// Package writer provides wrappers to the interface gxlog.Writer.
package writer

import "github.com/gxlog/gxlog"

// The Func type is a function wrapper to the interface gxlog.Writer.
// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(bs []byte, record *gxlog.Record)

// Write calls the underlying function. It implements the gxlog.Writer.
func (self Func) Write(bs []byte, record *gxlog.Record) {
	self(bs, record)
}
