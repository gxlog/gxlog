package syslog

import (
	"os"
	"path/filepath"

	"github.com/gxlog/gxlog"
)

type Facility int

// Facility definitions here to be cross compilation friendly
const (
	FacKern Facility = iota << 3
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
)

type Severity int

// Severity definitions here to be cross compilation friendly
const (
	SevEmerg Severity = iota
	SevAlert
	SevCrit
	SevErr
	SevWarning
	SevNotice
	SevInfo
	SevDebug
)

type Config struct {
	Tag         string
	Facility    Facility
	Network     string
	Addr        string
	SeverityMap map[gxlog.Level]Severity
	ReportOnErr bool
}

func NewConfig(tag string) *Config {
	if tag == "" {
		tag = filepath.Base(os.Args[0])
	}
	return &Config{
		Tag:         tag,
		Facility:    FacUser,
		ReportOnErr: true,
	}
}

func (this *Config) WithFacility(facility Facility) *Config {
	this.Facility = facility
	return this
}

func (this *Config) WithAddr(network, addr string) *Config {
	this.Network, this.Addr = network, addr
	return this
}

func (this *Config) WithSeverityMap(severityMap map[gxlog.Level]Severity) *Config {
	this.SeverityMap = severityMap
	return this
}

func (this *Config) WithReportOnErr(ok bool) *Config {
	this.ReportOnErr = ok
	return this
}
