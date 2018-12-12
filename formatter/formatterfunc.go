package formatter

import "github.com/gxlog/gxlog"

type FormatterFunc func(record *gxlog.Record) []byte

func (self FormatterFunc) Format(record *gxlog.Record) []byte {
	return self(record)
}
