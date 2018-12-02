package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

const DefaultTimeLayout = "2006-01-02 15:04:05.000000"

type timeFormatter struct {
	property string
	fmtspec  string
}

func newTimeFormatter(property, fmtspec string) *timeFormatter {
	if property == "" {
		property = DefaultTimeLayout
	}
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &timeFormatter{property: property, fmtspec: fmtspec}
}

func (this *timeFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := record.Time.Format(this.property)
	if this.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, desc)...)
	}
}
