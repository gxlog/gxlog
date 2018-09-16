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
		fmtspec = "%-8s"
	}
	return &levelFormatter{property: property, fmtspec: fmtspec}
}

func (this *levelFormatter) formatElement(record *gxlog.Record) string {
	var level string
	switch record.Level {
	case gxlog.LevelTrace:
		level = "TRACE"
	case gxlog.LevelDebug:
		level = "DEBUG"
	case gxlog.LevelInfo:
		level = "INFO"
	case gxlog.LevelNotice:
		level = "NOTICE"
	case gxlog.LevelWarning:
		level = "WARNING"
	case gxlog.LevelError:
		level = "ERROR"
	case gxlog.LevelCritical:
		level = "CRITICAL"
	case gxlog.LevelFatal:
		level = "FATAL"
	}
	return fmt.Sprintf(this.fmtspec, level)
}
