package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gratonos/gxlog"
)

type Writer struct {
	config Config

	file      *os.File
	pathname  string
	checkTime time.Time
	day       int
	fileSize  int64

	lock sync.Mutex
}

func Open(config *Config) (*Writer, error) {
	if config == nil {
		panic("nil config")
	}
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("file.Open: %v", err)
	}
	return &Writer{config: *config}, nil
}

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.closeFile(); err != nil {
		return fmt.Errorf("file.Close: %v", err)
	}
	return nil
}

func (this *Writer) Sync() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.file != nil {
		if err := this.file.Sync(); err != nil {
			return fmt.Errorf("file.Sync: %v", err)
		}
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.checkFile(record); err == nil {
		n, _ := this.file.Write(bs)
		this.fileSize += int64(n)
	}
}

func (this *Writer) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	return &copyConfig
}

func (this *Writer) SetConfig(config *Config) error {
	if config == nil {
		panic("nil config")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.setConfig(config); err != nil {
		return fmt.Errorf("file.SetConfig: %v", err)
	}
	return nil
}

func (this *Writer) UpdateConfig(fn func(*Config)) error {
	if fn == nil {
		panic("nil fn")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	fn(&copyConfig)
	if err := this.setConfig(&copyConfig); err != nil {
		return fmt.Errorf("file.UpdateConfig: %v", err)
	}
	return nil
}

func (this *Writer) checkFile(record *gxlog.Record) error {
	if this.file == nil ||
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

	this.file = file
	this.pathname = pathname
	this.day = record.Time.YearDay()
	this.fileSize = 0

	return nil
}

func (this *Writer) closeFile() error {
	if this.file != nil {
		if err := this.file.Close(); err != nil {
			return err
		}
		this.file = nil
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
	case DateStyleDash:
		fmtstr = "%04d-%02d-%02d"
	case DateStyleUnderscore:
		fmtstr = "%04d_%02d_%02d"
	case DateStyleDot:
		fmtstr = "%04d.%02d.%02d"
	}
	return fmt.Sprintf(fmtstr, tm.Year(), tm.Month(), tm.Day())
}

func (this *Writer) formatTime(tm time.Time) string {
	fmtstr := "%02d%02d%02d.%06d"
	switch this.config.TimeStyle {
	case TimeStyleDash:
		fmtstr = "%02d-%02d-%02d-%06d"
	case TimeStyleUnderscore:
		fmtstr = "%02d_%02d_%02d_%06d"
	case TimeStyleDot:
		fmtstr = "%02d.%02d.%02d.%06d"
	case TimeStyleColon:
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
