package text

import (
	"fmt"
	"strings"

	"github.com/gratonos/gxlog"
)

var gLevelDesc = []string{
	gxlog.LevelTrace: "TRACE",
	gxlog.LevelDebug: "DEBUG",
	gxlog.LevelInfo:  "INFO ",
	gxlog.LevelWarn:  "WARN ",
	gxlog.LevelError: "ERROR",
	gxlog.LevelFatal: "FATAL",
}

var gLevelDescChar = []string{
	gxlog.LevelTrace: "T",
	gxlog.LevelDebug: "D",
	gxlog.LevelInfo:  "I",
	gxlog.LevelWarn:  "W",
	gxlog.LevelError: "E",
	gxlog.LevelFatal: "F",
}

type levelFormatter struct {
	descList []string
	fmtspec  string
}

func newLevelFormatter(property, fmtspec string) *levelFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &levelFormatter{
		descList: selectDescList(property),
		fmtspec:  fmtspec,
	}
}

func (this *levelFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := this.descList[record.Level]
	if this.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, desc)...)
	}
}

func selectDescList(property string) []string {
	if strings.ToLower(property) == "char" {
		return gLevelDescChar
	}
	return gLevelDesc
}
