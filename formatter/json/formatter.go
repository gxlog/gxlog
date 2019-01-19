// Package json implements a json formatter which implements the Formatter.
package json

import (
	"strconv"
	"sync"
	"time"

	"github.com/gxlog/gxlog/formatter/internal/util"
	"github.com/gxlog/gxlog/iface"
)

// A Formatter implements the interface iface.Formatter.
//
// All methods of a Formatter are concurrency safe.
// A Formatter MUST be created with New.
type Formatter struct {
	config Config

	lock sync.Mutex
}

// New creates a new Formatter with the config.
func New(config Config) *Formatter {
	config.setDefaults()
	formatter := &Formatter{
		config: config,
	}
	return formatter
}

// Format implements the interface Formatter. It formats a Record.
func (formatter *Formatter) Format(record *iface.Record) []byte {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	buf := make([]byte, 0, formatter.config.MinBufSize)
	sep := ""
	buf = append(buf, "{"...)
	if formatter.config.Omit&Time == 0 {
		buf = formatStrField(buf, sep, "Time",
			record.Time.Format(time.RFC3339Nano), false)
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
	return append(buf, "}\n"...)
}

// Config returns the Config of the Formatter.
func (formatter *Formatter) Config() Config {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.config
}

// SetConfig sets the config to the Formatter.
func (formatter *Formatter) SetConfig(config Config) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	config.setDefaults()
	formatter.config = config
}

// UpdateConfig calls the fn with the Config of the Formatter, and then sets the
// returned Config to the Formatter. The fn must NOT be nil.
//
// Do NOT call any method of the Formatter or the Logger within the fn,
// or it may deadlock.
func (formatter *Formatter) UpdateConfig(fn func(Config) Config) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.config = fn(formatter.config)
}

func (formatter *Formatter) formatAux(buf []byte, sep string,
	aux *iface.Auxiliary) []byte {
	if formatter.config.Omit&Aux == Aux {
		return buf
	}
	if formatter.config.OmitEmpty&Aux == Aux &&
		aux.Prefix == "" && len(aux.Contexts) == 0 && !aux.Marked {
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
		!(formatter.config.OmitEmpty&Mark != 0 && !aux.Marked) {
		buf = formatBoolField(buf, sep, "Marked", aux.Marked)
	}
	return append(buf, "}"...)
}

func formatContexts(buf []byte, sep string, contexts []iface.Context) []byte {
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
	return append(buf, "]"...)
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
	return append(buf, `"`...)
}

func formatIntField(buf []byte, sep, key string, value int) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	return strconv.AppendInt(buf, int64(value), 10)
}

func formatBoolField(buf []byte, sep, key string, value bool) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	if value {
		return append(buf, "true"...)
	}
	return append(buf, "false"...)
}
