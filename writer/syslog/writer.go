// +build !nacl,!plan9,!windows

// Package syslog implements a syslog writer which implements the gxlog.Writer.
package syslog

import (
	"fmt"
	"log"
	"log/syslog"
	"sync"

	"github.com/gxlog/gxlog"
)

const severityMask = 0x07

type syslogFunc func(string) error

// A Writer implements the interface gxlog.Writer.
//
// All methods of a Writer are concurrency safe.
//
// A Writer must be created with Open.
type Writer struct {
	reportOnErr bool

	logFuncs [gxlog.LevelCount]syslogFunc
	writer   *syslog.Writer

	lock sync.Mutex
}

// Open creates a new Writer with the cfg. The cfg must NOT be nil.
func Open(cfg *Config) (*Writer, error) {
	wt, err := syslog.Dial(cfg.Network, cfg.Addr, syslog.Priority(cfg.Facility), cfg.Tag)
	if err != nil {
		return nil, fmt.Errorf("writer/syslog.Open: %v", err)
	}
	writer := &Writer{
		reportOnErr: cfg.ReportOnErr,
		writer:      wt,
	}
	severityMap := map[gxlog.Level]Severity{
		gxlog.Trace: SevDebug,
		gxlog.Debug: SevDebug,
		gxlog.Info:  SevInfo,
		gxlog.Warn:  SevWarning,
		gxlog.Error: SevErr,
		gxlog.Fatal: SevCrit,
	}
	writer.updateLogFuncs(severityMap)
	writer.updateLogFuncs(cfg.SeverityMap)
	return writer, nil
}

// Close closes the Writer.
func (writer *Writer) Close() error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.writer.Close(); err != nil {
		return fmt.Errorf("writer/syslog.Close: %v", err)
	}
	return nil
}

// Write implements the interface gxlog.Writer. It writes logs to the syslog.
func (writer *Writer) Write(bs []byte, record *gxlog.Record) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	err := writer.logFuncs[record.Level](string(bs))
	if writer.reportOnErr && err != nil {
		log.Println("writer/syslog.Write:", err)
	}
}

// ReportOnErr returns the reportOnErr of the Writer.
func (writer *Writer) ReportOnErr() bool {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.reportOnErr
}

// SetReportOnErr sets the reportOnErr of the Writer.
func (writer *Writer) SetReportOnErr(ok bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.reportOnErr = ok
}

// MapSeverity maps the severity of levels in the Writer by the severityMap.
// The severity of a level is left to be unchanged if it is not in the map.
func (writer *Writer) MapSeverity(severityMap map[gxlog.Level]Severity) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.updateLogFuncs(severityMap)
}

func (writer *Writer) updateLogFuncs(severityMap map[gxlog.Level]Severity) {
	for level, severity := range severityMap {
		var fn syslogFunc
		switch severity & severityMask {
		case SevDebug:
			fn = writer.writer.Debug
		case SevInfo:
			fn = writer.writer.Info
		case SevNotice:
			fn = writer.writer.Notice
		case SevWarning:
			fn = writer.writer.Warning
		case SevErr:
			fn = writer.writer.Err
		case SevCrit:
			fn = writer.writer.Crit
		case SevAlert:
			fn = writer.writer.Alert
		case SevEmerg:
			fn = writer.writer.Emerg
		}
		writer.logFuncs[level] = fn
	}
}
