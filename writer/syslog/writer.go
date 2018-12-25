// +build !nacl,!plan9,!windows

package syslog

import (
	"fmt"
	"log"
	"log/syslog"
	"sync"

	"github.com/gxlog/gxlog"
)

const cSeverityMask = 0x07

type syslogFunc func(string) error

type Writer struct {
	reportOnErr bool

	logFuncs [gxlog.LevelCount]syslogFunc
	writer   *syslog.Writer

	lock sync.Mutex
}

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

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.writer.Close(); err != nil {
		return fmt.Errorf("writer/syslog.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.lock.Lock()
	defer this.lock.Unlock()

	err := this.logFuncs[record.Level](string(bs))
	if this.reportOnErr && err != nil {
		log.Println("writer/syslog.Write:", err)
	}
}

func (this *Writer) ReportOnErr() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.reportOnErr
}

func (this *Writer) SetReportOnErr(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.reportOnErr = ok
}

func (this *Writer) MapSeverity(severityMap map[gxlog.Level]Severity) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.updateLogFuncs(severityMap)
}

func (this *Writer) updateLogFuncs(severityMap map[gxlog.Level]Severity) {
	for level, severity := range severityMap {
		var fn syslogFunc
		switch severity & cSeverityMask {
		case SevDebug:
			fn = this.writer.Debug
		case SevInfo:
			fn = this.writer.Info
		case SevNotice:
			fn = this.writer.Notice
		case SevWarning:
			fn = this.writer.Warning
		case SevErr:
			fn = this.writer.Err
		case SevCrit:
			fn = this.writer.Crit
		case SevAlert:
			fn = this.writer.Alert
		case SevEmerg:
			fn = this.writer.Emerg
		}
		this.logFuncs[level] = fn
	}
}
