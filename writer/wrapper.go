package writer

import (
	"io"

	"github.com/gratonos/gxlog"
)

type wrapper struct {
	writer io.Writer
}

func Wrap(wt io.Writer) gxlog.Writer {
	if wt == nil {
		panic("nil wt")
	}
	return &wrapper{writer: wt}
}

func (this *wrapper) Write(bs []byte, _ *gxlog.Record) {
	this.writer.Write(bs)
}
