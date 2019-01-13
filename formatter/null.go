package formatter

import (
	"github.com/gxlog/gxlog/iface"
)

var nullFormatter = Func(func(*iface.Record) []byte { return nil })

// Null returns the null Formatter.
func Null() iface.Formatter {
	return nullFormatter
}
