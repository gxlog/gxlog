package writer

import (
	"io"

	"github.com/gxlog/gxlog/iface"
)

type wrapper struct {
	writer io.Writer
}

// Wrap wraps a writer of io.Writer to iface.Writer. The writer must NOT be nil.
func Wrap(writer io.Writer) iface.Writer {
	if writer == nil {
		panic("writer.Wrap: nil writer")
	}
	return &wrapper{writer: writer}
}

func (wrap *wrapper) Write(bs []byte, _ *iface.Record) {
	wrap.writer.Write(bs)
}
