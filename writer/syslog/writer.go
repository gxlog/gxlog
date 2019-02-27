// Package syslog implements a syslog writer which implements the Writer.
//
// For performance and security, connect to the local syslog server and configure
// the local syslog server for log transmission if it is possible.
package syslog

import (
	"fmt"
	"sync"

	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/writer"
)

// A Writer implements the interface iface.Writer.
//
// All methods of a Writer are concurrency safe.
// A Writer MUST be created with Open.
type Writer struct {
	facility     Facility
	tag          string
	errorHandler writer.ErrorHandler

	severities []Severity
	log        *syslog

	lock sync.Mutex
}

// Open creates a new Writer with the config. If the Network field of the config
// is not specified, it will connect to the local syslog server with unix domain
// socket.
func Open(config Config) (*Writer, error) {
	config.setDefaults()
	log, err := syslogDial(config.Network, config.Addr)
	if err != nil {
		return nil, fmt.Errorf("writer/syslog.Open: %v", err)
	}
	severities := []Severity{
		iface.Trace: SevDebug,
		iface.Debug: SevDebug,
		iface.Info:  SevInfo,
		iface.Warn:  SevWarning,
		iface.Error: SevErr,
		iface.Fatal: SevCrit,
	}
	writer := &Writer{
		facility:     config.Facility,
		tag:          config.Tag,
		errorHandler: config.ErrorHandler,
		severities:   severities,
		log:          log,
	}
	writer.MapSeverities(config.SeverityMap)
	return writer, nil
}

// Close closes the Writer.
func (writer *Writer) Close() error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.log.Close(); err != nil {
		return fmt.Errorf("writer/syslog.Close: %v", err)
	}
	return nil
}

// Write implements the interface Writer. It writes logs to the syslog.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	severity := writer.severities[record.Level]
	priority := int(writer.facility) | int(severity)
	err := writer.log.Write(record.Time, priority, writer.tag, bs)
	if err != nil {
		writer.log.Close()
		if writer.errorHandler != nil {
			writer.errorHandler(bs, record, err)
		}
	}
}

// Facility returns the facility of the Writer.
func (writer *Writer) Facility() Facility {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.facility
}

// SetFacility sets the facility of the Writer.
func (writer *Writer) SetFacility(facility Facility) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.facility = facility
}

// Tag returns the tag of the Writer.
func (writer *Writer) Tag() string {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.tag
}

// SetTag sets the tag of the Writer.
func (writer *Writer) SetTag(tag string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.tag = tag
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

// MapSeverities maps the severity of levels according to the severityMap.
// The severity of a level is left to be unchanged if it is not in the map.
func (writer *Writer) MapSeverities(severityMap map[iface.Level]Severity) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	for level, severity := range severityMap {
		writer.severities[level] = severity
	}
}
