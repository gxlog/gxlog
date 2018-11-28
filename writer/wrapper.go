package writer

import (
	"io"

	"github.com/gratonos/gxlog"
)

type Wrapper struct {
	writer io.Writer
}

func (this *Wrapper) Write(bs []byte, _ *gxlog.Record) {
	this.writer.Write(bs)
}

func Wrap(wt io.Writer) *Wrapper {
	return &Wrapper{writer: wt}
}
