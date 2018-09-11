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

func (this *funcFormatter) formatElement(record *gxlog.Record) string {
	return fmt.Sprintf(this.fmtspec, record.Func)
}
