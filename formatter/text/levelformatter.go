package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type levelFormatter struct {
	property string
	fmtspec  string
}

func newLevelFormatter(property, fmtspec string) *levelFormatter {
	if fmtspec == "" {
		fmtspec = "%-5s"
	}
	return &levelFormatter{property: property, fmtspec: fmtspec}
}

func (this *levelFormatter) formatElement(record *gxlog.Record) string {
	var level string
	switch record.Level {
	case gxlog.LevelDebug:
		level = "DEBUG"
	case gxlog.LevelInfo:
		level = "INFO"
	case gxlog.LevelWarn:
		level = "WARN"
	case gxlog.LevelError:
		level = "ERROR"
	case gxlog.LevelFatal:
		level = "FATAL"
	default:
		level = "?????"
	}
	return fmt.Sprintf(this.fmtspec, level)
}
