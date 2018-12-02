package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

var levelDesc = []string{
	gxlog.LevelTrace: "TRACE",
	gxlog.LevelDebug: "DEBUG",
	gxlog.LevelInfo:  "INFO ",
	gxlog.LevelWarn:  "WARN ",
	gxlog.LevelError: "ERROR",
	gxlog.LevelFatal: "FATAL",
}

type levelFormatter struct {
	property string
	fmtspec  string
}

func newLevelFormatter(property, fmtspec string) *levelFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &levelFormatter{property: property, fmtspec: fmtspec}
}

func (this *levelFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := levelDesc[record.Level]
	if this.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, desc)...)
	}
}
