package syslog

import (
	"os"
	"path/filepath"

	"github.com/gxlog/gxlog"
)

type Priority int

// Severity definitions here to be cross compilation friendly
const (
	SevEmerg Priority = iota
	SevAlert
	SevCrit
	SevErr
	SevWarning
	SevNotice
	SevInfo
	SevDebug
)

// Facility definitions here to be cross compilation friendly
const (
	FacKern Priority = iota << 3
	FacUser
	FacMail
	FacDaemon
	FacAuth
	FacSyslog
	FacLPR
	FacNews
	FacUUCP
	FacCron
	FacAuthPriv
	FacFTP

	FacLocal0
	FacLocal1
	FacLocal2
	FacLocal3
	FacLocal4
	FacLocal5
	FacLocal6
	FacLocal7
)

const (
	DefaultFacility    = FacUser
	DefaultSeverity    = SevDebug
	DefaultPriority    = DefaultFacility | DefaultSeverity
	DefaultReportOnErr = true
)

const (
	DefaultTraceSeverity = SevDebug
	DefaultDebugSeverity = SevDebug
	DefaultInfoSeverity  = SevInfo
	DefaultWarnSeverity  = SevWarning
	DefaultErrorSeverity = SevErr
	DefaultFatalSeverity = SevCrit
)

type Config struct {
	Tag         string
	Priority    Priority
	Network     string
	Addr        string
	SeverityMap map[gxlog.Level]Priority
	ReportOnErr bool
}

func NewConfig(tag string) *Config {
	if tag == "" {
		tag = filepath.Base(os.Args[0])
	}
	return &Config{
		Tag:         tag,
		Priority:    DefaultPriority,
		ReportOnErr: DefaultReportOnErr,
	}
}

func (this *Config) WithPriority(priority Priority) *Config {
	this.Priority = priority
	return this
}

func (this *Config) WithAddr(network, addr string) *Config {
	this.Network, this.Addr = network, addr
	return this
}

func (this *Config) WithSeverityMap(severityMap map[gxlog.Level]Priority) *Config {
	this.SeverityMap = severityMap
	return this
}

func (this *Config) WithReportOnErr(ok bool) *Config {
	this.ReportOnErr = ok
	return this
}
