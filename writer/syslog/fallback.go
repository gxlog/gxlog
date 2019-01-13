// +build nacl plan9 windows

package syslog

import (
	"fmt"

	"github.com/gxlog/gxlog/iface"
)

const generalError = "not implemented on nacl, plan9 or windows"

type Writer struct{}

func Open(config Config) (*Writer, error) {
	return nil, fmt.Errorf("writer/syslog.Open: %s", generalError)
}

func (writer *Writer) Close() error {
	return fmt.Errorf("writer/syslog.Close: %s", generalError)
}

func (writer *Writer) Write(bs []byte, record *iface.Record) {}

func (writer *Writer) ReportOnErr() bool { return false }

func (writer *Writer) SetReportOnErr(ok bool) {}

func (writer *Writer) MapSeverity(severityMap map[iface.Level]Severity) {}
