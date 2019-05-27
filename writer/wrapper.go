package writer

import (
	"io"

	"github.com/gxlog/gxlog/iface"
)

type Wrapper struct {
	writer  io.Writer
	handler ErrorHandler
}

// Wrap wraps an io.Writer to iface.Writer. The writer must NOT be nil.
func Wrap(writer io.Writer, handler ErrorHandler) iface.Writer {
	return &Wrapper{
		writer:  writer,
		handler: handler,
	}
}

func (wrapper *Wrapper) Write(bs []byte, record *iface.Record) {
	_, err := wrapper.writer.Write(bs)
	if err != nil && wrapper.handler != nil {
		wrapper.handler(bs, record, err)
	}
}
