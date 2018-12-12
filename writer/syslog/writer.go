package syslog

import (
	"fmt"
	"log"
	"log/syslog"
	"sync"

	"github.com/gxlog/gxlog"
)

const (
	DefaultTraceSeverity = syslog.LOG_DEBUG
	DefaultDebugSeverity = syslog.LOG_DEBUG
	DefaultInfoSeverity  = syslog.LOG_INFO
	DefaultWarnSeverity  = syslog.LOG_WARNING
	DefaultErrorSeverity = syslog.LOG_ERR
	DefaultFatalSeverity = syslog.LOG_CRIT
)

const cSeverityMask = 0x07

type syslogFunc func(string) error

type Writer struct {
	reportOnErr bool

	logFuncs [gxlog.LevelCount]syslogFunc
	writer   *syslog.Writer

	lock sync.Mutex
}

func Open(config *Config) (*Writer, error) {
	if config == nil {
		panic("nil config")
	}
	wt, err := syslog.Dial(config.Network, config.Addr, config.Priority, config.Tag)
	if err != nil {
		return nil, fmt.Errorf("syslog.Open: %v", err)
	}
	writer := &Writer{
		reportOnErr: config.ReportOnErr,
		writer:      wt,
	}
	severityMap := map[gxlog.Level]syslog.Priority{
		gxlog.LevelTrace: DefaultTraceSeverity,
		gxlog.LevelDebug: DefaultDebugSeverity,
		gxlog.LevelInfo:  DefaultInfoSeverity,
		gxlog.LevelWarn:  DefaultWarnSeverity,
		gxlog.LevelError: DefaultErrorSeverity,
		gxlog.LevelFatal: DefaultFatalSeverity,
	}
	writer.updateLogFuncs(severityMap)
	writer.updateLogFuncs(config.SeverityMap)
	return writer, nil
}

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.writer.Close(); err != nil {
		return fmt.Errorf("syslog.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.lock.Lock()
	defer this.lock.Unlock()

	err := this.logFuncs[record.Level](string(bs))
	if this.reportOnErr && err != nil {
		log.Println("syslog.Write:", err)
	}
}

func (this *Writer) MapSeverity(severityMap map[gxlog.Level]syslog.Priority) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.updateLogFuncs(severityMap)
}

func (this *Writer) updateLogFuncs(severityMap map[gxlog.Level]syslog.Priority) {
	for level, severity := range severityMap {
		var fn syslogFunc
		switch severity & cSeverityMask {
		case syslog.LOG_DEBUG:
			fn = this.writer.Debug
		case syslog.LOG_INFO:
			fn = this.writer.Info
		case syslog.LOG_NOTICE:
			fn = this.writer.Notice
		case syslog.LOG_WARNING:
			fn = this.writer.Warning
		case syslog.LOG_ERR:
			fn = this.writer.Err
		case syslog.LOG_CRIT:
			fn = this.writer.Crit
		case syslog.LOG_ALERT:
			fn = this.writer.Alert
		case syslog.LOG_EMERG:
			fn = this.writer.Emerg
		}
		this.logFuncs[level] = fn
	}
}
