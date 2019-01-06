// Package formatter provides wrappers of the interface iface.Formatter.
package formatter

import "github.com/gxlog/gxlog/iface"

// The Func type is a function wrapper of the interface iface.Formatter.
// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(record *iface.Record) []byte

// Format calls the underlying function. It implements the iface.Formatter.
func (fn Func) Format(record *iface.Record) []byte {
	return fn(record)
}
