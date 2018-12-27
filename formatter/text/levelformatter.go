package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog"
)

var levelDesc = []string{
	gxlog.Trace: "TRACE",
	gxlog.Debug: "DEBUG",
	gxlog.Info:  "INFO ",
	gxlog.Warn:  "WARN ",
	gxlog.Error: "ERROR",
	gxlog.Fatal: "FATAL",
}

var levelDescChar = []string{
	gxlog.Trace: "T",
	gxlog.Debug: "D",
	gxlog.Info:  "I",
	gxlog.Warn:  "W",
	gxlog.Error: "E",
	gxlog.Fatal: "F",
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

func (formatter *levelFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := formatter.descList[record.Level]
	if formatter.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(formatter.fmtspec, desc)...)
	}
}

func selectDescList(property string) []string {
	if strings.ToLower(property) == "char" {
		return levelDescChar
	}
	return levelDesc
}
