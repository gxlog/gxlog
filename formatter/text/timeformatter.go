package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog"
)

const (
	cFmtDate  = "2006-01-02"
	cFmtTime  = "15:04:05"
	cFmtMilli = ".000"
	cFmtMicro = ".000000"
	cFmtNano  = ".000000000"
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

func (this *timeFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	desc := record.Time.Format(this.layout)
	if this.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, desc)...)
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
		layout = cFmtDate + " " + cFmtTime
	case "time":
		layout = cFmtTime
	default:
		return "2006-01-02 15:04:05.000000"
	}
	switch decimalType {
	case "ms":
		layout += cFmtMilli
	case "us":
		layout += cFmtMicro
	case "ns":
		layout += cFmtNano
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
