// Package writer provides wrappers to the interface Writer.
package writer

import (
	"github.com/gxlog/gxlog/iface"
)

// The Func type is a function wrapper to the interface Writer.
// Do NOT call any method of the Logger within the function, or it may deadlock.
type Func func(bs []byte, record *iface.Record)

// Write calls the underlying function. It implements the interface Writer.
func (fn Func) Write(bs []byte, record *iface.Record) {
	fn(bs, record)
}
