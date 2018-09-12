package null

import "github.com/gratonos/gxlog"

type Writer struct{}

func (Writer) Write([]byte, *gxlog.Record) {}
