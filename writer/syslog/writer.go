// +build !nacl,!plan9,!windows

// Package syslog implements a syslog writer which implements the Writer.
package syslog

import (
	"fmt"
	"log/syslog"
	"sync"

	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/writer"
)

const severityMask = 0x07

type syslogFunc func(string) error

// A Writer implements the interface iface.Writer.
//
// All methods of a Writer are concurrency safe.
// A Writer MUST be created with Open.
type Writer struct {
	errorHandler writer.ErrorHandler

	logFuncs [iface.LevelCount + iface.Trace]syslogFunc
	writer   *syslog.Writer

	lock sync.Mutex
}

// Open creates a new Writer with the config.
func Open(config Config) (*Writer, error) {
	config.setDefaults()
	wt, err := syslog.Dial(config.Network, config.Addr,
		syslog.Priority(config.Facility), config.Tag)
	if err != nil {
		return nil, fmt.Errorf("writer/syslog.Open: %v", err)
	}
	writer := &Writer{
		errorHandler: config.ErrorHandler,
		writer:       wt,
	}
	severityMap := map[iface.Level]Severity{
		iface.Trace: SevDebug,
		iface.Debug: SevDebug,
		iface.Info:  SevInfo,
		iface.Warn:  SevWarning,
		iface.Error: SevErr,
		iface.Fatal: SevCrit,
	}
	writer.updateLogFuncs(severityMap)
	writer.updateLogFuncs(config.SeverityMap)
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

// Write implements the interface Writer. It writes logs to the syslog.
//
// NOTICE: the standard syslog package will get the timestamp itself which is a
// tiny bit later than Record.Time.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	err := writer.logFuncs[record.Level](string(bs))
	if err != nil && writer.errorHandler != nil {
		writer.errorHandler(bs, record, err)
	}
}

// ErrorHandler returns the error handler of the Writer.
func (writer *Writer) ErrorHandler() writer.ErrorHandler {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.errorHandler
}

// SetErrorHandler sets the error handler of the Writer.
func (writer *Writer) SetErrorHandler(handler writer.ErrorHandler) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.errorHandler = handler
}

// MapSeverity maps the severity of levels according to the severityMap.
// The severity of a level is left to be unchanged if it is not in the map.
func (writer *Writer) MapSeverity(severityMap map[iface.Level]Severity) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.updateLogFuncs(severityMap)
}

func (writer *Writer) updateLogFuncs(severityMap map[iface.Level]Severity) {
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
