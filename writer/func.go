// Package writer provides wrappers to the interface iface.Writer.
package writer

import (
	"github.com/gxlog/gxlog/iface"
)

// The Func type is a function wrapper to the interface iface.Writer.
// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(bs []byte, record *iface.Record)

// Write calls the underlying function. It implements the iface.Writer.
func (fn Func) Write(bs []byte, record *iface.Record) {
	fn(bs, record)
}
