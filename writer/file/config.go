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
	DefaultExt           = ".log"
	DefaultSeparator     = "."
	DefaultDateStyle     = DateStyleCompact
	DefaultTimeStyle     = TimeStyleCompact
	DefaultMaxFileSize   = 20 * 1024 * 1024
	DefaultCheckInterval = time.Second * 5
	DefaultNewDirEachDay = true
)

type Config struct {
	Path          string
	Base          string
	Ext           string
	Separator     string
	DateStyle     DateStyleID
	TimeStyle     TimeStyleID
	MaxFileSize   int64
	CheckInterval time.Duration
	NewDirEachDay bool
}

func NewConfig(path, base string) *Config {
	return &Config{
		Path:          path,
		Base:          base,
		Ext:           DefaultExt,
		Separator:     DefaultSeparator,
		DateStyle:     DefaultDateStyle,
		TimeStyle:     DefaultTimeStyle,
		MaxFileSize:   DefaultMaxFileSize,
		CheckInterval: DefaultCheckInterval,
		NewDirEachDay: DefaultNewDirEachDay,
	}
}
