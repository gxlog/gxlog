package text

import (
	"fmt"
	"strconv"

	"github.com/gxlog/gxlog"
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

func (formatter *lineFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if formatter.fmtspec == "%d" {
		return append(buf, strconv.Itoa(record.Line)...)
	} else {
		return append(buf, fmt.Sprintf(formatter.fmtspec, record.Line)...)
	}
}
