package writer

import (
	"github.com/gxlog/gxlog/iface"
)

// Multi creates a Writer that duplicates its writes to all the writers.
// Any writer must NOT be nil.
func Multi(writers ...iface.Writer) iface.Writer {
	return Func(func(bs []byte, record *iface.Record) {
		for _, writer := range writers {
			writer.Write(bs, record)
		}
	})
}
