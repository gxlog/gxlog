package writer

import (
	"github.com/gxlog/gxlog/iface"
)

var nullWriter = Func(func([]byte, *iface.Record) {})

// Null returns the null Writer.
func Null() iface.Writer {
	return nullWriter
}
