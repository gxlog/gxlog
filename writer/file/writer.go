package file

import (
	"compress/flate"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gxlog/gxlog"
)

type Writer struct {
	config Config

	writer    io.WriteCloser
	pathname  string
	checkTime time.Time
	day       int
	fileSize  int64

	lock sync.Mutex
}

func Open(config *Config) (*Writer, error) {
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("writer/file.Open: %v", err)
	}
	return &Writer{config: *config}, nil
}

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.closeFile(); err != nil {
		return fmt.Errorf("writer/file.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.lock.Lock()
	defer this.lock.Unlock()

	err := this.checkFile(record)
	if err == nil {
		var n int
		n, err = this.writer.Write(bs)
		this.fileSize += int64(n)
	}
	if this.config.ReportOnErr && err != nil {
		log.Println("writer/file.Write:", err)
	}
}

func (this *Writer) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	return &copyConfig
}

func (this *Writer) SetConfig(config *Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.setConfig(config); err != nil {
		return fmt.Errorf("writer/file.SetConfig: %v", err)
	}
	return nil
}

func (this *Writer) UpdateConfig(fn func(Config) Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	config := fn(this.config)
	if err := this.setConfig(&config); err != nil {
		return fmt.Errorf("writer/file.UpdateConfig: %v", err)
	}
	return nil
}

func (this *Writer) checkFile(record *gxlog.Record) error {
	if this.writer == nil ||
		this.day != record.Time.YearDay() ||
		this.fileSize >= this.config.MaxFileSize {
		return this.createFile(record)
	} else if time.Since(this.checkTime) >= this.config.CheckInterval {
		this.checkTime = time.Now()
		if _, err := os.Stat(this.pathname); err != nil {
			return this.createFile(record)
		}
	}
	return nil
}

func (this *Writer) createFile(record *gxlog.Record) error {
	if err := this.closeFile(); err != nil {
		return err
	}

	path := this.formatPath(record.Time)
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	filename := this.formatFilename(record.Time)
	pathname := filepath.Join(path, filename)
	file, err := os.Create(pathname)
	if err != nil {
		return err
	}

	var writer io.WriteCloser = file
	if this.config.AESKey != "" {
		// newAESWriter will return the input writer when an error occurs
		writer, err = newAESWriter(writer, this.config.AESKey, this.config.BlockMode)
		if err != nil {
			writer.Close()
			return err
		}
	}
	if this.config.GzipLevel != flate.NoCompression {
		// newGzipWriter will return the input writer when an error occurs
		writer, err = newGzipWriter(writer, this.config.GzipLevel)
		if err != nil {
			writer.Close()
			return err
		}
	}

	this.writer = writer
	this.pathname = pathname
	this.day = record.Time.YearDay()
	this.fileSize = 0

	return nil
}

func (this *Writer) closeFile() error {
	if this.writer != nil {
		if err := this.writer.Close(); err != nil {
			return err
		}
		this.writer = nil
	}
	return nil
}

func (this *Writer) formatPath(tm time.Time) string {
	path := this.config.Path
	if this.config.NewDirEachDay {
		path = filepath.Join(path, this.formatDate(tm))
	}
	return path
}

func (this *Writer) formatFilename(tm time.Time) string {
	elements := []string{}
	if this.config.Base != "" {
		elements = append(elements, this.config.Base)
	}
	if !this.config.NewDirEachDay {
		elements = append(elements, this.formatDate(tm))
	}
	elements = append(elements, this.formatTime(tm))
	return strings.Join(elements, this.config.Separator) + this.config.Ext
}

func (this *Writer) formatDate(tm time.Time) string {
	fmtstr := "%04d%02d%02d"
	switch this.config.DateStyle {
	case DateDash:
		fmtstr = "%04d-%02d-%02d"
	case DateUnderscore:
		fmtstr = "%04d_%02d_%02d"
	case DateDot:
		fmtstr = "%04d.%02d.%02d"
	}
	return fmt.Sprintf(fmtstr, tm.Year(), tm.Month(), tm.Day())
}

func (this *Writer) formatTime(tm time.Time) string {
	fmtstr := "%02d%02d%02d.%06d"
	switch this.config.TimeStyle {
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

func (this *Writer) needNewFile(config *Config) bool {
	if config.Path != this.config.Path ||
		config.Base != this.config.Base ||
		config.Ext != this.config.Ext ||
		config.Separator != this.config.Separator ||
		config.DateStyle != this.config.DateStyle ||
		config.TimeStyle != this.config.TimeStyle ||
		config.GzipLevel != this.config.GzipLevel ||
		config.AESKey != this.config.AESKey ||
		config.BlockMode != this.config.BlockMode ||
		config.NewDirEachDay != this.config.NewDirEachDay {
		return true
	}
	return false
}

func (this *Writer) setConfig(config *Config) error {
	if err := config.Check(); err != nil {
		return err
	}
	if this.needNewFile(config) {
		if err := this.closeFile(); err != nil {
			return err
		}
	}
	this.config = *config
	return nil
}
