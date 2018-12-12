// +build nacl plan9 windows

package syslog

import (
	"fmt"

	"github.com/gxlog/gxlog"
)

const cError = "not implemented on nacl, plan9 or windows"

type Writer struct{}

func Open(cfg *Config) (*Writer, error) {
	return nil, fmt.Errorf("syslog.Open: %s", cError)
}

func (this *Writer) Close() error {
	return fmt.Errorf("syslog.Close: %s", cError)
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {}

func (this *Writer) MapSeverity(severityMap map[gxlog.Level]Priority) {}
