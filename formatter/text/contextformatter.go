package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type contextFormatter struct {
	property string
	fmtspec  string
	buf      []byte
}

func newContextFormatter(property, fmtspec string) *contextFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &contextFormatter{property: property, fmtspec: fmtspec}
}

func (this *contextFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return this.format(buf, record.Contexts)
	} else {
		this.buf = this.buf[:0]
		this.buf = this.format(this.buf, record.Contexts)
		return append(buf, fmt.Sprintf(this.fmtspec, this.buf)...)
	}
}

func (this *contextFormatter) format(buf []byte, contexts []gxlog.Context) []byte {
	if len(contexts) != 0 {
		buf = append(buf, '[')
	}
	for _, ctx := range contexts {
		buf = append(buf, '(')
		buf = append(buf, ctx.Key...)
		buf = append(buf, ':')
		buf = append(buf, ctx.Value...)
		buf = append(buf, ')')
	}
	if len(contexts) != 0 {
		buf = append(buf, ']')
	}
	return buf
}
