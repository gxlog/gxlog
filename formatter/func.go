package formatter

import "github.com/gxlog/gxlog"

type Func func(record *gxlog.Record) []byte

func (self Func) Format(record *gxlog.Record) []byte {
	return self(record)
}
