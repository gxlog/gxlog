package file

import (
	"errors"
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

func (this *Config) WithExt(ext string) *Config {
	this.Ext = ext
	return this
}

func (this *Config) WithSeparator(sep string) *Config {
	this.Separator = sep
	return this
}

func (this *Config) WithDateStyle(style DateStyleID) *Config {
	this.DateStyle = style
	return this
}

func (this *Config) WithTimeStyle(style TimeStyleID) *Config {
	this.TimeStyle = style
	return this
}

func (this *Config) WithMaxFileSize(size int64) *Config {
	this.MaxFileSize = size
	return this
}

func (this *Config) WithCheckInterval(interval time.Duration) *Config {
	this.CheckInterval = interval
	return this
}

func (this *Config) WithNewDirEachDay(ok bool) *Config {
	this.NewDirEachDay = ok
	return this
}

func (this *Config) Check() error {
	if this.MaxFileSize <= 0 {
		return errors.New("Config.MaxFileSize must be greater than 0")
	}
	if this.CheckInterval <= 0 {
		return errors.New("Config.CheckInterval must be greater than 0")
	}
	return nil
}
