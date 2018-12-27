// Package formatter provides wrappers of the interface gxlog.Formatter.
package formatter

import "github.com/gxlog/gxlog"

// The Func type is a function wrapper of the interface gxlog.Formatter.
// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(record *gxlog.Record) []byte

// Format calls the underlying function. It implements the gxlog.Formatter.
func (self Func) Format(record *gxlog.Record) []byte {
	return self(record)
}
