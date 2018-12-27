package writer

import (
	"io"

	"github.com/gxlog/gxlog"
)

type wrapper struct {
	writer io.Writer
}

// Wrap wraps a writer of io.Writer to gxlog.Writer. The writer must NOT be nil.
func Wrap(writer io.Writer) gxlog.Writer {
	if writer == nil {
		panic("writer.Wrap: nil writer")
	}
	return &wrapper{writer: writer}
}

func (wrap *wrapper) Write(bs []byte, _ *gxlog.Record) {
	wrap.writer.Write(bs)
}
