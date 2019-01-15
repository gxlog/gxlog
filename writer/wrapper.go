package writer

import (
	"io"

	"github.com/gxlog/gxlog/iface"
)

// Wrap wraps an io.Writer to iface.Writer. The writer must NOT be nil.
func Wrap(writer io.Writer) iface.Writer {
	return Func(func(bs []byte, _ *iface.Record) {
		writer.Write(bs)
	})
}
