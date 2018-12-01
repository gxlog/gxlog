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

func (this *contextFormatter) FormatElement(record *gxlog.Record) string {
	this.buf = this.buf[:0]

	if len(record.Contexts) != 0 {
		this.buf = append(this.buf, '[')
	}
	for _, ctx := range record.Contexts {
		this.buf = append(this.buf, '(')
		this.buf = append(this.buf, ctx.Key...)
		this.buf = append(this.buf, ':')
		this.buf = append(this.buf, ctx.Value...)
		this.buf = append(this.buf, ')')
	}
	if len(record.Contexts) != 0 {
		this.buf = append(this.buf, ']')
	}

	if this.fmtspec == "%s" {
		return string(this.buf)
	} else {
		return fmt.Sprintf(this.fmtspec, this.buf)
	}
}
