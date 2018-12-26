package formatter

import "github.com/gxlog/gxlog"

// Do NOT call methods of the Logger within the function, or it will deadlock.
type Func func(record *gxlog.Record) []byte

func (self Func) Format(record *gxlog.Record) []byte {
	return self(record)
}
