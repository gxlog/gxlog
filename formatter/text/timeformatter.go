package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog"
)

const (
	dateLayout    = "2006-01-02"
	timeLayout    = "15:04:05"
	milliLayout   = ".000"
	microLayout   = ".000000"
	nanoLayout    = ".000000000"
	defaultLayout = "2006-01-02 15:04:05.000000"
)

type timeFormatter struct {
	layout  string
	fmtspec string
}

func newTimeFormatter(property, fmtspec string) *timeFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &timeFormatter{
		layout:  makeTimeLayout(property),
		fmtspec: fmtspec,
	}
}

func (formatter *timeFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := record.Time.Format(formatter.layout)
	if formatter.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(formatter.fmtspec, desc)...)
	}
}

func makeTimeLayout(property string) string {
	if strings.ContainsAny(property, "0123456789") {
		return property
	}

	var layout string
	timeType, decimalType := getTimeOptions(property)
	switch timeType {
	case "date":
		layout = dateLayout + " " + timeLayout
	case "time":
		layout = timeLayout
	default:
		return defaultLayout
	}
	switch decimalType {
	case "ms":
		layout += milliLayout
	case "us":
		layout += microLayout
	case "ns":
		layout += nanoLayout
	}
	return layout
}

func getTimeOptions(str string) (string, string) {
	fields := strings.Split(strings.ToLower(str), ".")
	if len(fields) == 0 {
		return "", ""
	}
	if len(fields) == 1 {
		return fields[0], ""
	}
	return fields[0], fields[1]
}
