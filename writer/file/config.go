package file

import (
	"compress/flate"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gxlog/gxlog/writer"
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
type Config struct {
	// Path is the location to where logs output.
	// Shell expansion is NOT supported.
	// When it is modified in a file writer, a new log file will be created.
	// If Path is not specified, "." is used.
	Path string
	// Base is the first segment of the name of log files.
	// When it is modified in a file writer, a new log file will be created.
	// If Base is not specified, filepath.Base(os.Args[0]).<pid> is used.
	Base string
	// Ext is the extension name of log files.
	// When it is modified in a file writer, a new log file will be created.
	// If Ext is not specified, ".log" is used.
	Ext string
	// Separator is the segment separator of name of log files.
	// When it is modified in a file writer, a new log file will be created.
	// If Separator is not specified, "." is used.
	Separator string
	// DateStyle is the date format style for naming log files.
	// When it is modified in a file writer, a new log file will be created.
	// If DateStyle is not specified, DateCompact is used.
	DateStyle DateStyle
	// TimeStyle is the time format style for naming log files.
	// When it is modified in a file writer, a new log file will be created.
	// If TimeStyle is not specified, TimeCompact is used.
	TimeStyle TimeStyle
	// MaxFileSize is the max size of a log file BEFORE compression because
	// (*gzip.Writer).Write returns the count of bytes before compression.
	// If MaxFileSize is not specified, (20 * 1024 * 1024) is used.
	// It must NOT be negative.
	MaxFileSize int64
	// CheckInterval is the time interval to check whether the current log file
	// still exists. If not, a new log file will be created.
	// It is useful when you want to remove all log files and do not want to
	// restart the process.
	// If CheckInterval is not specified, (time.Second * 5) is used.
	// For performance, it is better NOT to be less than 1s.
	CheckInterval time.Duration
	// GzipLevel is the level of gzip of log files. It will be handled by package
	// compress/gzip. It MUST be flate.DefaultCompression, flate.NoCompression,
	// flate.HuffmanOnly or any integer value between flate.BestSpeed and
	// flate.BestCompression inclusive.
	// When it is modified in a file writer, a new log file will be created.
	// If GzipLevel is not specified, flate.NoCompression is used.
	GzipLevel int
	// AESKey is a hexadecimal encoded AES key. It MUST be either empty, 128 bits,
	// 192 bits or 256 bits, e.g. 70856575b161fbcca8fc12e1f70fc1c8.
	// If it is not empty, the AES encryption is enabled. Each log file will have
	// an independent initialization vector.
	// When it is modified in a file writer, a new log file will be created.
	AESKey string
	// BlockMode is the block mode of AES. It MUST be either CFB, CTR or OFB.
	// When it is modified in a file writer, a new log file will be created.
	// If BlockMode is not specified, CFB is used.
	BlockMode BlockCipherMode
	// NoDirForDays specifies NOT to create a new directory each day.
	// If NoDirForDays is true, the pattern of name of log files is
	// <base><sep><date><sep><time><ext>, otherwise it is <base><sep><time><ext>.
	// When it is modified in a file writer, a new log file will be created.
	NoDirForDays bool
	// ErrorHandler will be called when an error occurs if it is not nil.
	ErrorHandler writer.ErrorHandler
}

func (config *Config) setDefaults() {
	if config.Path == "" {
		config.Path = "."
	}
	if config.Base == "" {
		config.Base = filepath.Base(os.Args[0]) + "." + strconv.Itoa(os.Getpid())
	}
	if config.Ext == "" {
		config.Ext = ".log"
	}
	if config.Separator == "" {
		config.Separator = "."
	}
	if config.MaxFileSize == 0 {
		config.MaxFileSize = 20 * 1024 * 1024
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = time.Second * 5
	}
}

func (config *Config) check() error {
	if config.MaxFileSize < 0 {
		return errors.New("Config.MaxFileSize must NOT be negative")
	}
	if config.CheckInterval < 0 {
		return errors.New("Config.CheckInterval must NOT be negative")
	}
	if config.GzipLevel < flate.HuffmanOnly ||
		config.GzipLevel > flate.BestCompression {
		return errors.New("Config.GzipLevel is invalid")
	}
	key, err := hex.DecodeString(config.AESKey)
	if err != nil {
		return errors.New("Config.AESKey is invalid")
	}
	keyLen := len(key)
	if keyLen != 0 && keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return errors.New("Config.AESKey is invalid")
	}
	return nil
}
