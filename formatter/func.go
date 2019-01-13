// Package formatter provides wrappers to the interface Formatter.
package formatter

import (
	"github.com/gxlog/gxlog/iface"
)

// The Func type is a function wrapper to the interface Formatter.
// Do NOT call any method of the Logger within the function, or it may deadlock.
type Func func(record *iface.Record) []byte

// Format calls the underlying function. It implements the interface Formatter.
func (fn Func) Format(record *iface.Record) []byte {
	return fn(record)
}
