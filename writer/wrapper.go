package writer

import (
	"io"

	"github.com/gratonos/gxlog"
)

type wrapper struct {
	writer io.Writer
}

func Wrap(wt io.Writer) gxlog.Writer {
	return &wrapper{writer: wt}
}

func (this *wrapper) Write(bs []byte, _ *gxlog.Record) {
	this.writer.Write(bs)
}
