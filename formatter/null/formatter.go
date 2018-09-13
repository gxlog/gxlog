package null

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter"
)

var Formatter = formatter.FormatterFunc(func(*gxlog.Record) []byte { return nil })
