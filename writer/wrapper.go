package writer

import (
	"io"

	"github.com/gxlog/gxlog"
)

type wrapper struct {
	writer io.Writer
}

func Wrap(writer io.Writer) gxlog.Writer {
	if writer == nil {
		panic("writer.Wrap: nil writer")
	}
	return &wrapper{writer: writer}
}

func (this *wrapper) Write(bs []byte, _ *gxlog.Record) {
	this.writer.Write(bs)
}
