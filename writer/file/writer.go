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

func (this *Writer) Start(config *Config) {
	this.config = *config
}

func (this *Writer) Stop() {
	this.closeFile()
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	if err := this.checkFile(record); err == nil {
		n, _ := this.file.Write(bs)
		this.fileSize += int64(n)
	}
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
	this.closeFile()

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

func (this *Writer) closeFile() {
	if this.file != nil {
		this.file.Close()
		this.file = nil
	}
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
	milli := fmt.Sprintf("%09d", tm.Nanosecond())[:3]
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
	return strings.Join([]string{hour, minute, second}, sep) + "." + milli
}
