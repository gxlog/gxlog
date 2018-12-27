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

func (cfg *Config) WithExt(ext string) *Config {
	cfg.Ext = ext
	return cfg
}

func (cfg *Config) WithSeparator(sep string) *Config {
	cfg.Separator = sep
	return cfg
}

func (cfg *Config) WithDateStyle(style DateStyleID) *Config {
	cfg.DateStyle = style
	return cfg
}

func (cfg *Config) WithTimeStyle(style TimeStyleID) *Config {
	cfg.TimeStyle = style
	return cfg
}

func (cfg *Config) WithMaxFileSize(size int64) *Config {
	cfg.MaxFileSize = size
	return cfg
}

func (cfg *Config) WithCheckInterval(interval time.Duration) *Config {
	cfg.CheckInterval = interval
	return cfg
}

func (cfg *Config) WithGzipLevel(level int) *Config {
	cfg.GzipLevel = level
	return cfg
}

func (cfg *Config) WithAESKey(key string) *Config {
	cfg.AESKey = key
	return cfg
}

func (cfg *Config) WithBlockMode(mode BlockCipherMode) *Config {
	cfg.BlockMode = mode
	return cfg
}

func (cfg *Config) WithNewDirEachDay(ok bool) *Config {
	cfg.NewDirEachDay = ok
	return cfg
}

func (cfg *Config) WithReportOnErr(ok bool) *Config {
	cfg.ReportOnErr = ok
	return cfg
}

func (cfg *Config) Check() error {
	if cfg.MaxFileSize <= 0 {
		return errors.New("Config.MaxFileSize must be greater than 0")
	}
	if cfg.CheckInterval <= 0 {
		return errors.New("Config.CheckInterval must be greater than 0")
	}
	if cfg.GzipLevel < flate.HuffmanOnly || cfg.GzipLevel > flate.BestCompression {
		return errors.New("Config.GzipLevel must be DefaultCompression, NoCompression, " +
			"HuffmanOnly or any integer value between BestSpeed and BestCompression inclusive")
	}
	key, err := hex.DecodeString(cfg.AESKey)
	if err != nil {
		return errors.New("Config.AESKey must be hexadecimal encoded without prefix 0X or 0x")
	}
	keyLen := len(key)
	if keyLen != 0 && keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return errors.New("Config.AESKey must be either empty, 128 bits, 192 bits or 256 bits")
	}
	if cfg.BlockMode < CFB || cfg.BlockMode > OFB {
		return errors.New("Config.BlockMode must be either CFB, CTR or OFB")
	}
	return nil
}
