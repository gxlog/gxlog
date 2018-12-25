package file

import (
	"compress/flate"
	"encoding/hex"
	"errors"
	"time"
)

type DateStyleID int

const (
	DateCompact DateStyleID = iota
	DateDash
	DateUnderscore
	DateDot
)

type TimeStyleID int

const (
	TimeCompact TimeStyleID = iota
	TimeDash
	TimeUnderscore
	TimeDot
	TimeColon
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
	GzipLevel     int
	AESKey        string
	BlockMode     BlockCipherMode
	NewDirEachDay bool
	ReportOnErr   bool
}

func NewConfig(path, base string) *Config {
	return &Config{
		Path:          path,
		Base:          base,
		Ext:           ".log",
		Separator:     ".",
		DateStyle:     DateCompact,
		TimeStyle:     TimeCompact,
		MaxFileSize:   20 * 1024 * 1024,
		CheckInterval: time.Second * 5,
		GzipLevel:     flate.NoCompression,
		BlockMode:     CFB,
		NewDirEachDay: true,
		ReportOnErr:   true,
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

func (this *Config) WithGzipLevel(level int) *Config {
	this.GzipLevel = level
	return this
}

func (this *Config) WithAESKey(key string) *Config {
	this.AESKey = key
	return this
}

func (this *Config) WithBlockMode(mode BlockCipherMode) *Config {
	this.BlockMode = mode
	return this
}

func (this *Config) WithNewDirEachDay(ok bool) *Config {
	this.NewDirEachDay = ok
	return this
}

func (this *Config) WithReportOnErr(ok bool) *Config {
	this.ReportOnErr = ok
	return this
}

func (this *Config) Check() error {
	if this.MaxFileSize <= 0 {
		return errors.New("Config.MaxFileSize must be greater than 0")
	}
	if this.CheckInterval <= 0 {
		return errors.New("Config.CheckInterval must be greater than 0")
	}
	if this.GzipLevel < flate.HuffmanOnly || this.GzipLevel > flate.BestCompression {
		return errors.New("Config.GzipLevel must be DefaultCompression, NoCompression, " +
			"HuffmanOnly or any integer value between BestSpeed and BestCompression inclusive")
	}
	key, err := hex.DecodeString(this.AESKey)
	if err != nil {
		return errors.New("Config.AESKey must be hexadecimal encoded without prefix 0X or 0x")
	}
	keyLen := len(key)
	if keyLen != 0 && keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return errors.New("Config.AESKey must be either empty, 128 bits, 192 bits or 256 bits")
	}
	if this.BlockMode < CFB || this.BlockMode > OFB {
		return errors.New("Config.BlockMode must be either CFB, CTR or OFB")
	}
	return nil
}
