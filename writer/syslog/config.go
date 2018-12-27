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

func (cfg *Config) WithFacility(facility Facility) *Config {
	cfg.Facility = facility
	return cfg
}

func (cfg *Config) WithAddr(network, addr string) *Config {
	cfg.Network, cfg.Addr = network, addr
	return cfg
}

func (cfg *Config) WithSeverityMap(severityMap map[gxlog.Level]Severity) *Config {
	cfg.SeverityMap = severityMap
	return cfg
}

func (cfg *Config) WithReportOnErr(ok bool) *Config {
	cfg.ReportOnErr = ok
	return cfg
}
