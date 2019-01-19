// Package file implements a file writer which implements the Writer.
package file

import (
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gxlog/gxlog/iface"
)

// A Writer implements the interface iface.Writer.
//
// All methods of a Writer are concurrency safe.
// A Writer MUST be created with Open.
type Writer struct {
	config Config

	writer    io.WriteCloser
	pathname  string
	checkTime time.Time
	day       int
	fileSize  int64

	lock sync.Mutex
}

// Open creates a new Writer with the config.
func Open(config Config) (*Writer, error) {
	config.setDefaults()
	if err := config.check(); err != nil {
		return nil, fmt.Errorf("writer/file.Open: %v", err)
	}
	return &Writer{config: config}, nil
}

// Close closes the Writer.
func (writer *Writer) Close() error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.closeFile(); err != nil {
		return fmt.Errorf("writer/file.Close: %v", err)
	}
	return nil
}

// Write implements the interface Writer. It writes logs to files.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	err := writer.checkFile(record)
	if err == nil {
		var n int
		n, err = writer.writer.Write(bs)
		writer.fileSize += int64(n)
	}
	if err != nil && writer.config.ErrorHandler != nil {
		writer.config.ErrorHandler(bs, record, err)
	}
}

// Config returns the Config of the Writer.
func (writer *Writer) Config() Config {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.config
}

// SetConfig sets the config to the Writer.
// If the config is invalid, it returns an error and the Config of the Writer
// is left to be unchanged.
func (writer *Writer) SetConfig(config Config) error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.setConfig(&config); err != nil {
		return fmt.Errorf("writer/file.SetConfig: %v", err)
	}
	return nil
}

// UpdateConfig calls the fn with the Config of the Writer, and then
// sets the returned config to the Writer. The fn must NOT be nil.
// If the returned config is invalid, it returns an error and the Config of
// the Writer is left to be unchanged.
//
// Do NOT call any method of the Writer or the Logger within the fn,
// or it may deadlock.
func (writer *Writer) UpdateConfig(fn func(Config) Config) error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	config := fn(writer.config)
	if err := writer.setConfig(&config); err != nil {
		return fmt.Errorf("writer/file.UpdateConfig: %v", err)
	}
	return nil
}

func (writer *Writer) checkFile(record *iface.Record) error {
	if writer.writer == nil ||
		writer.day != record.Time.YearDay() ||
		writer.fileSize >= writer.config.MaxFileSize {
		return writer.createFile(record)
	} else if time.Since(writer.checkTime) >= writer.config.CheckInterval {
		writer.checkTime = time.Now()
		if _, err := os.Stat(writer.pathname); err != nil {
			return writer.createFile(record)
		}
	}
	return nil
}

func (writer *Writer) createFile(record *iface.Record) error {
	if err := writer.closeFile(); err != nil {
		return err
	}

	path := writer.formatPath(record.Time)
	if err := os.MkdirAll(path, writer.config.DirPerm); err != nil {
		return err
	}

	filename := writer.formatFilename(record.Time)
	pathname := filepath.Join(path, filename)
	file, err := os.Create(pathname)
	if err != nil {
		return err
	}

	var wt io.WriteCloser = file
	if writer.config.AESKey != "" {
		// newAESWriter will return the input writer when an error occurs
		wt, err = newAESWriter(wt, writer.config.AESKey, writer.config.BlockMode)
		if err != nil {
			wt.Close()
			return err
		}
	}
	if writer.config.GzipLevel != flate.NoCompression {
		// newGzipWriter will return the input writer when an error occurs
		wt, err = newGzipWriter(wt, writer.config.GzipLevel)
		if err != nil {
			wt.Close()
			return err
		}
	}

	writer.writer = wt
	writer.pathname = pathname
	writer.day = record.Time.YearDay()
	writer.fileSize = 0

	return nil
}

func (writer *Writer) closeFile() error {
	if writer.writer != nil {
		if err := writer.writer.Close(); err != nil {
			return err
		}
		writer.writer = nil
	}
	return nil
}

func (writer *Writer) formatPath(tm time.Time) string {
	path := writer.config.Path
	if !writer.config.NoDirForDays {
		path = filepath.Join(path, writer.formatDate(tm))
	}
	return path
}

func (writer *Writer) formatFilename(tm time.Time) string {
	elements := []string{}
	if writer.config.Base != "" {
		elements = append(elements, writer.config.Base)
	}
	if writer.config.NoDirForDays {
		elements = append(elements, writer.formatDate(tm))
	}
	elements = append(elements, writer.formatTime(tm))
	return strings.Join(elements, writer.config.Separator) + writer.config.Ext
}

func (writer *Writer) formatDate(tm time.Time) string {
	fmtstr := "%04d%02d%02d"
	switch writer.config.DateStyle {
	case DateDash:
		fmtstr = "%04d-%02d-%02d"
	case DateUnderscore:
		fmtstr = "%04d_%02d_%02d"
	case DateDot:
		fmtstr = "%04d.%02d.%02d"
	}
	return fmt.Sprintf(fmtstr, tm.Year(), tm.Month(), tm.Day())
}

func (writer *Writer) formatTime(tm time.Time) string {
	fmtstr := "%02d%02d%02d.%06d"
	switch writer.config.TimeStyle {
	case TimeDash:
		fmtstr = "%02d-%02d-%02d-%06d"
	case TimeUnderscore:
		fmtstr = "%02d_%02d_%02d_%06d"
	case TimeDot:
		fmtstr = "%02d.%02d.%02d.%06d"
	case TimeColon:
		fmtstr = "%02d:%02d:%02d.%06d"
	}
	return fmt.Sprintf(fmtstr, tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond()/1000)
}

func (writer *Writer) needNewFile(config *Config) bool {
	if config.Path != writer.config.Path ||
		config.Base != writer.config.Base ||
		config.Ext != writer.config.Ext ||
		config.Separator != writer.config.Separator ||
		config.DateStyle != writer.config.DateStyle ||
		config.TimeStyle != writer.config.TimeStyle ||
		config.GzipLevel != writer.config.GzipLevel ||
		config.AESKey != writer.config.AESKey ||
		config.BlockMode != writer.config.BlockMode ||
		config.NoDirForDays != writer.config.NoDirForDays {
		return true
	}
	return false
}

func (writer *Writer) setConfig(config *Config) error {
	config.setDefaults()
	if err := config.check(); err != nil {
		return err
	}
	if writer.needNewFile(config) {
		if err := writer.closeFile(); err != nil {
			return err
		}
	}
	writer.config = *config
	return nil
}
