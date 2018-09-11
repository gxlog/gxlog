package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type msgFormatter struct {
	property string
	fmtspec  string
}

func newMsgFormatter(property, fmtspec string) *msgFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgFormatter{property: property, fmtspec: fmtspec}
}

func (this *msgFormatter) formatElement(record *gxlog.Record) string {
	return fmt.Sprintf(this.fmtspec, record.Msg)
}
