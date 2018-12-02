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

func (this *levelFormatter) FormatElement(record *gxlog.Record) string {
	if this.fmtspec == "%s" {
		return levelDesc[record.Level]
	} else {
		return fmt.Sprintf(this.fmtspec, levelDesc[record.Level])
	}
}
