package syslog

import (
	"os"
	"path/filepath"

	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/writer"
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
type Config struct {
	// If Tag is not specified, filepath.Base(os.Args[0]) is used.
	Tag string
	// If Facility is not specified, FacKern is used.
	Facility Facility
	// If Network is not specified, it will connect to the local syslog server
	// with unix domain socket. Otherwise, Network will be passed to net.Dial.
	Network string
	// Addr will be passed to net.Dial if Network is not empty.
	Addr string
	// SeverityMap is used to remap the severity of levels.
	// The severity of a level is left to be unchanged if it is not in the map.
	// The default mapping is as the follows:
	//   Trace: SevDebug
	//   Debug: SevDebug
	//   Info:  SevInfo
	//   Warn:  SevWarning
	//   Error: SevErr
	//   Fatal: SevCrit
	SeverityMap map[iface.Level]Severity
	// ErrorHandler will be called when an error occurs if it is not nil.
	ErrorHandler writer.ErrorHandler
}

func (config *Config) setDefaults() {
	if config.Tag == "" {
		config.Tag = filepath.Base(os.Args[0])
	}
}
