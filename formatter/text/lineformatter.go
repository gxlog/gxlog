package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog"
)

type lineFormatter struct {
	property string
	fmtspec  string
}

func newLineFormatter(property, fmtspec string) *lineFormatter {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineFormatter{property: property, fmtspec: fmtspec}
}

func (this *lineFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%d" {
		return append(buf, strconv.Itoa(record.Line)...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Line)...)
	}
}
