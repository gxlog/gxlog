package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog/iface"
)

var levelDesc = []string{
	iface.Trace: "TRACE",
	iface.Debug: "DEBUG",
	iface.Info:  "INFO ",
	iface.Warn:  "WARN ",
	iface.Error: "ERROR",
	iface.Fatal: "FATAL",
}

var levelDescChar = []string{
	iface.Trace: "T",
	iface.Debug: "D",
	iface.Info:  "I",
	iface.Warn:  "W",
	iface.Error: "E",
	iface.Fatal: "F",
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

func (formatter *levelFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	desc := formatter.descList[record.Level]
	if formatter.fmtspec == "%s" {
		return append(buf, desc...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, desc)...)
}

func selectDescList(property string) []string {
	if strings.ToLower(property) == "char" {
		return levelDescChar
	}
	return levelDesc
}
