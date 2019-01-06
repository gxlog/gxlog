package syslog

import (
	"os"
	"path/filepath"

	"github.com/gxlog/gxlog/iface"
)

// The Facility defines the facility type of syslog.
type Facility int

// Facility definitions here to be cross compilation friendly.
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

// The Severity defines the severity type of syslog.
type Severity int

// Severity definitions here to be cross compilation friendly.
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

// A Config is used to configure a syslog writer.
// A Config should be created with NewConfig.
type Config struct {
	// If Tag is empty, it will be filepath.Base(os.Args[0]).
	Tag      string
	Facility Facility
	// If Network is empty, it will connect to the local syslog server.
	// Otherwise, it will be passed to net.Dial.
	Network string
	// Addr will be passed to net.Dial if Network is not empty.
	Addr string
	// SeverityMap is used to remap the severity of each level.
	// The severity of a level is left to be unchanged if it is not in the map.
	// The default mapping is as the follows:
	//     iface.Trace: SevDebug
	//     iface.Debug: SevDebug
	//     iface.Info:  SevInfo
	//     iface.Warn:  SevWarning
	//     iface.Error: SevErr
	//     iface.Fatal: SevCrit
	SeverityMap map[iface.Level]Severity
	// ReportOnErr specifies whether to report errors by log.Println.
	ReportOnErr bool
}

// NewConfig creates a new Config. By default, the Facility is FacUser and
// the ReportOnErr is true.
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

// WithFacility sets the Facility of the Config and returns the Config.
func (cfg *Config) WithFacility(facility Facility) *Config {
	cfg.Facility = facility
	return cfg
}

// WithAddr sets the Network and Addr of the Config and returns the Config.
func (cfg *Config) WithAddr(network, addr string) *Config {
	cfg.Network, cfg.Addr = network, addr
	return cfg
}

// WithSeverityMap sets the SeverityMap of the Config and returns the Config.
func (cfg *Config) WithSeverityMap(severityMap map[iface.Level]Severity) *Config {
	cfg.SeverityMap = severityMap
	return cfg
}

// WithReportOnErr sets the ReportOnErr of the Config and returns the Config.
func (cfg *Config) WithReportOnErr(ok bool) *Config {
	cfg.ReportOnErr = ok
	return cfg
}
