// Package json implements a json formatter which implements the gxlog.Formatter.
package json

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/internal/util"
)

// A Formatter implements the interface gxlog.Formatter.
//
// All methods of a Formatter are concurrency safe.
//
// A Formatter must be created with New.
type Formatter struct {
	config Config

	lock sync.Mutex
}

// New creates a new Formatter with the config. The config must not be nil.
func New(config *Config) *Formatter {
	if config.MinBufSize < 0 {
		panic("formatter/json.New: Config.MinBufSize must not be negative")
	}
	formatter := &Formatter{
		config: *config,
	}
	return formatter
}

// Format implements the interface gxlog.Formatter. It formats a Record.
func (formatter *Formatter) Format(record *gxlog.Record) []byte {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	buf := make([]byte, 0, formatter.config.MinBufSize)
	sep := ""
	buf = append(buf, "{"...)
	if formatter.config.Omit&Time == 0 {
		buf = formatStrField(buf, sep, "Time", record.Time.Format(time.RFC3339Nano), false)
		sep = ","
	}
	if formatter.config.Omit&Level == 0 {
		buf = formatIntField(buf, sep, "Level", int(record.Level))
		sep = ","
	}
	if formatter.config.Omit&File == 0 {
		file := util.LastSegments(record.File, formatter.config.FileSegs, '/')
		buf = formatStrField(buf, sep, "File", file, true)
		sep = ","
	}
	if formatter.config.Omit&Line == 0 {
		buf = formatIntField(buf, sep, "Line", record.Line)
		sep = ","
	}
	if formatter.config.Omit&Pkg == 0 {
		pkg := util.LastSegments(record.Pkg, formatter.config.PkgSegs, '/')
		buf = formatStrField(buf, sep, "Pkg", pkg, false)
		sep = ","
	}
	if formatter.config.Omit&Func == 0 {
		fn := util.LastSegments(record.Func, formatter.config.FuncSegs, '.')
		buf = formatStrField(buf, sep, "Func", fn, false)
		sep = ","
	}
	if formatter.config.Omit&Msg == 0 {
		buf = formatStrField(buf, sep, "Msg", record.Msg, true)
		sep = ","
	}
	buf = formatter.formatAux(buf, sep, &record.Aux)
	buf = append(buf, "}\n"...)
	return buf
}

// Config returns a copy of config of the Formatter.
func (formatter *Formatter) Config() *Config {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	copyConfig := formatter.config
	return &copyConfig
}

// SetConfig sets the copy of config to the Formatter. The config must NOT be nil.
// If the config is invalid, it returns an error and the config of the Formatter
// is left to be unchanged.
func (formatter *Formatter) SetConfig(config *Config) error {
	if config.MinBufSize < 0 {
		return errors.New("formatter/json.SetConfig: Config.MinBufSize must not be negative")
	}

	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.config = *config
	return nil
}

// UpdateConfig will call fn with copy of the config of the Formatter, and then
// sets copy of the returned config to the Formatter. The fn must NOT be nil.
// If the returned config is invalid, it returns an error and the config of
// the Formatter is left to be unchanged.
//
// Do NOT call methods of the Formatter within fn, or it will deadlock.
func (formatter *Formatter) UpdateConfig(fn func(Config) Config) error {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	config := fn(formatter.config)

	if config.MinBufSize < 0 {
		return errors.New("formatter/json.UpdateConfig: Config.MinBufSize must not be negative")
	}
	formatter.config = config
	return nil
}

func (formatter *Formatter) formatAux(buf []byte, sep string, aux *gxlog.Auxiliary) []byte {
	if formatter.config.Omit&Aux == Aux {
		return buf
	}
	if formatter.config.OmitEmpty&Aux == Aux &&
		aux.Prefix == "" && len(aux.Contexts) == 0 && aux.Marked == false {
		return buf
	}
	buf = append(buf, sep...)
	sep = ""
	buf = append(buf, `"Aux":{`...)
	if formatter.config.Omit&Prefix == 0 &&
		!(formatter.config.OmitEmpty&Prefix != 0 && aux.Prefix == "") {
		buf = formatStrField(buf, sep, "Prefix", aux.Prefix, true)
		sep = ","
	}
	if formatter.config.Omit&Context == 0 &&
		!(formatter.config.OmitEmpty&Context != 0 && len(aux.Contexts) == 0) {
		buf = formatContexts(buf, sep, aux.Contexts)
		sep = ","
	}
	if formatter.config.Omit&Mark == 0 &&
		!(formatter.config.OmitEmpty&Mark != 0 && aux.Marked == false) {
		buf = formatBoolField(buf, sep, "Marked", aux.Marked)
	}
	buf = append(buf, "}"...)
	return buf
}

func formatContexts(buf []byte, sep string, contexts []gxlog.Context) []byte {
	buf = append(buf, sep...)
	sep = ""
	if len(contexts) == 0 {
		return append(buf, `"Contexts":null`...)
	}
	buf = append(buf, `"Contexts":[`...)
	for _, context := range contexts {
		buf = append(buf, sep...)
		buf = append(buf, "{"...)
		buf = formatStrField(buf, "", "Key", context.Key, true)
		buf = formatStrField(buf, ",", "Value", context.Value, true)
		buf = append(buf, "}"...)
		sep = ","
	}
	buf = append(buf, "]"...)
	return buf
}

func formatStrField(buf []byte, sep, key, value string, esc bool) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":"`...)
	if esc {
		buf = escape(buf, value)
	} else {
		buf = append(buf, value...)
	}
	buf = append(buf, `"`...)
	return buf
}

func formatIntField(buf []byte, sep, key string, value int) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	buf = append(buf, strconv.Itoa(value)...)
	return buf
}

func formatBoolField(buf []byte, sep, key string, value bool) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	if value {
		buf = append(buf, "true"...)
	} else {
		buf = append(buf, "false"...)
	}
	return buf
}
