package json

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/internal/util"
)

type Formatter struct {
	config Config

	lock sync.Mutex
}

func New(config *Config) *Formatter {
	if config.MinBufSize < 0 {
		panic("formatter/json.New: Config.MinBufSize must not be negative")
	}
	formatter := &Formatter{
		config: *config,
	}
	return formatter
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	this.lock.Lock()
	defer this.lock.Unlock()

	buf := make([]byte, 0, this.config.MinBufSize)
	sep := ""
	buf = append(buf, "{"...)
	if this.config.Omit&Time == 0 {
		buf = formatStrField(buf, sep, "Time", record.Time.Format(time.RFC3339Nano), false)
		sep = ","
	}
	if this.config.Omit&Level == 0 {
		buf = formatIntField(buf, sep, "Level", int(record.Level))
		sep = ","
	}
	if this.config.Omit&File == 0 {
		file := util.LastSegments(record.File, this.config.FileSegs, '/')
		buf = formatStrField(buf, sep, "File", file, true)
		sep = ","
	}
	if this.config.Omit&Line == 0 {
		buf = formatIntField(buf, sep, "Line", record.Line)
		sep = ","
	}
	if this.config.Omit&Pkg == 0 {
		pkg := util.LastSegments(record.Pkg, this.config.PkgSegs, '/')
		buf = formatStrField(buf, sep, "Pkg", pkg, false)
		sep = ","
	}
	if this.config.Omit&Func == 0 {
		fn := util.LastSegments(record.Func, this.config.FuncSegs, '.')
		buf = formatStrField(buf, sep, "Func", fn, false)
		sep = ","
	}
	if this.config.Omit&Msg == 0 {
		buf = formatStrField(buf, sep, "Msg", record.Msg, true)
		sep = ","
	}
	buf = this.formatAux(buf, sep, &record.Aux)
	buf = append(buf, "}\n"...)
	return buf
}

func (this *Formatter) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	return &copyConfig
}

func (this *Formatter) SetConfig(config *Config) error {
	if config.MinBufSize < 0 {
		return errors.New("formatter/json.SetConfig: Config.MinBufSize must not be negative")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	this.config = *config
	return nil
}

func (this *Formatter) UpdateConfig(fn func(Config) Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	config := fn(this.config)

	if config.MinBufSize < 0 {
		return errors.New("formatter/json.UpdateConfig: Config.MinBufSize must not be negative")
	}
	this.config = config
	return nil
}

func (this *Formatter) formatAux(buf []byte, sep string, aux *gxlog.Auxiliary) []byte {
	if this.config.Omit&Aux == Aux {
		return buf
	}
	if this.config.OmitEmpty&Aux == Aux &&
		aux.Prefix == "" && len(aux.Contexts) == 0 && aux.Marked == false {
		return buf
	}
	buf = append(buf, sep...)
	sep = ""
	buf = append(buf, `"Aux":{`...)
	if this.config.Omit&Prefix == 0 &&
		!(this.config.OmitEmpty&Prefix != 0 && aux.Prefix == "") {
		buf = formatStrField(buf, sep, "Prefix", aux.Prefix, true)
		sep = ","
	}
	if this.config.Omit&Context == 0 &&
		!(this.config.OmitEmpty&Context != 0 && len(aux.Contexts) == 0) {
		buf = formatContexts(buf, sep, aux.Contexts)
		sep = ","
	}
	if this.config.Omit&Mark == 0 &&
		!(this.config.OmitEmpty&Mark != 0 && aux.Marked == false) {
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
