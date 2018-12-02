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

func (this *lineFormatter) FormatElement(record *gxlog.Record) string {
	if this.fmtspec == "%d" {
		return strconv.Itoa(record.Line)
	} else {
		return fmt.Sprintf(this.fmtspec, record.Line)
	}
}
