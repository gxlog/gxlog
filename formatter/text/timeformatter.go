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

func (this *timeFormatter) formatElement(record *gxlog.Record) string {
	desc := record.Time.Format(this.property)
	if this.fmtspec == "%s" {
		return desc
	} else {
		return fmt.Sprintf(this.fmtspec, desc)
	}
}
