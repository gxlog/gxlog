package syslog

import (
	"log/syslog"
	"os"
	"path/filepath"

	"github.com/gxlog/gxlog"
)

const (
	DefaultFacility    = syslog.LOG_USER
	DefaultSeverity    = syslog.LOG_DEBUG
	DefaultPriority    = DefaultFacility | DefaultSeverity
	DefaultReportOnErr = true
)

type Config struct {
	Tag         string
	Priority    syslog.Priority
	Network     string
	Addr        string
	SeverityMap map[gxlog.Level]syslog.Priority
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

func (this *Config) WithPriority(priority syslog.Priority) *Config {
	this.Priority = priority
	return this
}

func (this *Config) WithAddr(network, addr string) *Config {
	this.Network, this.Addr = network, addr
	return this
}

func (this *Config) WithSeverityMap(severityMap map[gxlog.Level]syslog.Priority) *Config {
	this.SeverityMap = severityMap
	return this
}

func (this *Config) WithReportOnErr(ok bool) *Config {
	this.ReportOnErr = ok
	return this
}
