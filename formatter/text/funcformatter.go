package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type funcFormatter struct {
	property string
	fmtspec  string
}

func newFuncFormatter(property, fmtspec string) *funcFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &funcFormatter{property: property, fmtspec: fmtspec}
}

func (this *funcFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Func...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Func)...)
	}
}
