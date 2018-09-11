package text

import (
	"fmt"

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

func (this *lineFormatter) formatElement(record *gxlog.Record) string {
	return fmt.Sprintf(this.fmtspec, record.Line)
}
