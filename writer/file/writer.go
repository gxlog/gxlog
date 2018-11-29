package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gratonos/gxlog"
)

type Writer struct {
	config    Config
	file      *os.File
	pathname  string
	checkTime time.Time
	day       int
	fileSize  int64
}

func Open(config *Config) (*Writer, error) {
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("file.Open: %v", err)
	}
	return &Writer{config: *config}, nil
}

func (this *Writer) Close() error {
	if err := this.closeFile(); err != nil {
		return fmt.Errorf("file.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	if err := this.checkFile(record); err == nil {
		n, _ := this.file.Write(bs)
		this.fileSize += int64(n)
	}
}

func (this *Writer) GetConfig() Config {
	return this.config
}

func (this *Writer) SetConfig(config *Config) error {
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

func (this *Writer) UpdateConfig(fn func(*Config)) error {
	config := this.config
	fn(&config)
	return this.SetConfig(&config)
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
	year := fmt.Sprintf("%04d", tm.Year())
	month := fmt.Sprintf("%02d", tm.Month())
	day := fmt.Sprintf("%02d", tm.Day())
	sep := ""
	switch this.config.DateStyle {
	case DateStyleDash:
		sep = "-"
	case DateStyleUnderscore:
		sep = "_"
	case DateStyleDot:
		sep = "."
	}
	return strings.Join([]string{year, month, day}, sep)
}

func (this *Writer) formatTime(tm time.Time) string {
	hour := fmt.Sprintf("%02d", tm.Hour())
	minute := fmt.Sprintf("%02d", tm.Minute())
	second := fmt.Sprintf("%02d", tm.Second())
	micro := fmt.Sprintf("%09d", tm.Nanosecond())[:6]
	sep := ""
	switch this.config.TimeStyle {
	case TimeStyleDash:
		sep = "-"
	case TimeStyleUnderscore:
		sep = "_"
	case TimeStyleDot:
		sep = "."
	case TimeStyleColon:
		sep = ":"
	}
	return strings.Join([]string{hour, minute, second}, sep) + "." + micro
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
