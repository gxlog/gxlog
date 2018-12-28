package file

import (
	"compress/flate"
	"encoding/hex"
	"errors"
	"time"
)

// The DateStyle defines the type of date format style for naming log files.
type DateStyle int

// All available date format styles here.
const (
	// YYYYMMDD
	DateCompact DateStyle = iota
	// YYYY-MM-DD
	DateDash
	// YYYY_MM_DD
	DateUnderscore
	// YYYY.MM.DD
	DateDot
)

// The TimeStyle defines the type of time format style for naming log files.
type TimeStyle int

// All available time format styles here.
const (
	// hhmmss.uuuuuu
	TimeCompact TimeStyle = iota
	// hh-mm-ss-uuuuuu
	TimeDash
	// hh_mm_ss_uuuuuu
	TimeUnderscore
	// hh.mm.ss.uuuuuu
	TimeDot
	// hh:mm:ss.uuuuuu
	TimeColon
)

// A Config is used to configure a file writer.
// A Config should be created with NewConfig.
type Config struct {
	// Path is the path to where logs output.
	Path string
	// Base is the first segment of log files' name. It is better NOT to be empty.
	Base string
	// Ext is the extension of log files' name.
	Ext string
	// Separator is the segment separator of log files' name.
	Separator string
	// DateStyle is the date format style for naming log files.
	DateStyle DateStyle
	// TimeStyle is the time format style for naming log files.
	TimeStyle TimeStyle
	// MaxFileSize is the max size of a log file BEFORE compression because of
	// that (*gzip.Writer).Write returns the count of bytes before compression.
	// It must be positive.
	MaxFileSize int64
	// CheckInterval is the interval to check whether the current log file still
	// exists. If not, a new log file will be created.
	// It is useful when you want to remove all log files and do not want to
	// restart the process. For performance, it is better NOT to be less than 1s.
	CheckInterval time.Duration
	// GzipLevel is the level of gzip of log files. It will be handled by package
	// compress/gzip. It must be flate.DefaultCompression, flate.NoCompression,
	// flate.HuffmanOnly or any integer value between flate.BestSpeed and
	// flate.BestCompression inclusive.
	GzipLevel int
	// AESKey is a hexadecimal encoded AES key. It must be either empty, 128 bits,
	// 192 bits or 256 bits, e.g. 70856575b161fbcca8fc12e1f70fc1c8.
	// If it is not empty, the AES encryption is enabled.
	AESKey string
	// Available block modes are CFB, CTR and OFB.
	BlockMode BlockCipherMode
	// NewDirEachDay specifies whether to create a new directory each day.
	// If NewDirEachDay is true, the pattern of log files' name is
	// <base><sep><time><ext>, otherwise it is <base><sep><date><sep><time><ext>.
	NewDirEachDay bool
	// ReportOnErr specifies whether to report errors by log.Println.
	ReportOnErr bool
}

// NewConfig creates a new Config. It is better that the base is NOT empty.
// By default, the Ext is ".log", Separator is ".", DateStyle is DateCompact,
// TimeStyle is TimeCompact, MaxFileSize is 20MB, CheckInterval is 5s,
// GzipLevel is flate.NoCompression, AESKey is empty, BlockMode is CFB,
// NewDirEachDay is true, ReportOnErr is true.
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

// WithExt sets the Ext of the Config and returns the Config.
func (cfg *Config) WithExt(ext string) *Config {
	cfg.Ext = ext
	return cfg
}

// WithSeparator sets the Separator of the Config and returns the Config.
func (cfg *Config) WithSeparator(sep string) *Config {
	cfg.Separator = sep
	return cfg
}

// WithDateStyle sets the DateStyle of the Config and returns the Config.
func (cfg *Config) WithDateStyle(style DateStyle) *Config {
	cfg.DateStyle = style
	return cfg
}

// WithTimeStyle sets the TimeStyle of the Config and returns the Config.
func (cfg *Config) WithTimeStyle(style TimeStyle) *Config {
	cfg.TimeStyle = style
	return cfg
}

// WithMaxFileSize sets the MaxFileSize of the Config and returns the Config.
func (cfg *Config) WithMaxFileSize(size int64) *Config {
	cfg.MaxFileSize = size
	return cfg
}

// WithCheckInterval sets the CheckInterval of the Config and returns the Config.
func (cfg *Config) WithCheckInterval(interval time.Duration) *Config {
	cfg.CheckInterval = interval
	return cfg
}

// WithGzipLevel sets the GzipLevel of the Config and returns the Config.
func (cfg *Config) WithGzipLevel(level int) *Config {
	cfg.GzipLevel = level
	return cfg
}

// WithAESKey sets the AESKey of the Config and returns the Config.
func (cfg *Config) WithAESKey(key string) *Config {
	cfg.AESKey = key
	return cfg
}

// WithBlockMode sets the BlockMode of the Config and returns the Config.
func (cfg *Config) WithBlockMode(mode BlockCipherMode) *Config {
	cfg.BlockMode = mode
	return cfg
}

// WithNewDirEachDay sets the NewDirEachDay of the Config and returns the Config.
func (cfg *Config) WithNewDirEachDay(ok bool) *Config {
	cfg.NewDirEachDay = ok
	return cfg
}

// WithReportOnErr sets the ReportOnErr of the Config and returns the Config.
func (cfg *Config) WithReportOnErr(ok bool) *Config {
	cfg.ReportOnErr = ok
	return cfg
}

// Check returns an error if the Config is invalid, otherwise it returns nil.
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
