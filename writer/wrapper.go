package writer

import (
	"io"

	"github.com/gratonos/gxlog"
)

type Wrapper struct {
	wt io.Writer
}

func (this *Wrapper) Write(bs []byte, _ *gxlog.Record) {
	this.wt.Write(bs)
}

func Wrap(wt io.Writer) *Wrapper {
	return &Wrapper{wt: wt}
}
