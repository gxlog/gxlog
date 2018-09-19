package file

import (
	"time"
)

type DateStyleID int

const (
	DateStyleCompact DateStyleID = iota
	DateStyleDash
	DateStyleUnderscore
	DateStyleDot
)

type TimeStyleID int

const (
	TimeStyleCompact TimeStyleID = iota
	TimeStyleDash
	TimeStyleUnderscore
	TimeStyleDot
	TimeStyleColon
)

const (
	DefaultPath          = "log"
	DefaultBase          = ""
	DefaultExt           = ".log"
	DefaultDateStyle     = DateStyleCompact
	DefaultTimeStyle     = TimeStyleCompact
	DefaultSeparator     = "."
	DefaultMaxFileSize   = 20 * 1024 * 1024
	DefaultCheckInterval = time.Second * 5
	DefaultNewDirEachDay = true
)

type Config struct {
	Path          string
	Base          string
	Ext           string
	DateStyle     DateStyleID
	TimeStyle     TimeStyleID
	Separator     string
	MaxFileSize   int64
	CheckInterval time.Duration
	NewDirEachDay bool
}

func NewConfig() *Config {
	return &Config{
		Path:          DefaultPath,
		Base:          DefaultBase,
		Ext:           DefaultExt,
		DateStyle:     DefaultDateStyle,
		TimeStyle:     DefaultTimeStyle,
		Separator:     DefaultSeparator,
		MaxFileSize:   DefaultMaxFileSize,
		CheckInterval: DefaultCheckInterval,
		NewDirEachDay: DefaultNewDirEachDay,
	}
}
