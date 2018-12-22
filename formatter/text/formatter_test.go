package text_test

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/text"
)

const (
	cLayout      = "2006-01-02 15:04:05.000000000"
	cDate        = "2018-08-01"
	cTime        = "07:12:07"
	cDecimal     = "235605270"
	cLevel       = gxlog.LevelInfo
	cFile        = "/home/test/data/src/go/workspace/src/github.com/gxlog/gxlog/logger.go"
	cLine        = 64
	cPkg         = "github.com/gxlog/gxlog"
	cFunc        = "Test"
	cMsg         = "testing"
	cPrefix      = "**** "
	cContextPair = "(k1: v1) (k2: v2)"
	cContextList = "k1: v1, k2: v2"
)

var gTmplContexts = []gxlog.Context{
	{"k1", "v1"},
	{"k2", "v2"},
}

var gTmplRecord gxlog.Record

func init() {
	clock, err := time.ParseInLocation(cLayout, cDate+" "+cTime+"."+cDecimal, time.Local)
	if err != nil {
		panic(err)
	}

	gTmplRecord = gxlog.Record{
		Time:  clock,
		Level: cLevel,
		File:  cFile,
		Line:  cLine,
		Pkg:   cPkg,
		Func:  cFunc,
		Msg:   cMsg,
		Aux: gxlog.Auxiliary{
			Prefix:   cPrefix,
			Contexts: gTmplContexts,
			Marked:   true,
		},
	}
}

func TestDefaultHeader(t *testing.T) {
	formatter := text.New(text.NewConfig())
	expect := fmt.Sprintf("%s %s.%s %s %s:%d %s.%s %s[%s] %s\n",
		cDate, cTime, cDecimal[:6], "INFO ", cFile, cLine, cPkg, cFunc,
		cPrefix, cContextPair, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)
}

func TestCompactHeader(t *testing.T) {
	formatter := text.New(text.NewConfig().WithHeader(text.CompactHeader))
	expect := fmt.Sprintf("%s.%s %s %s:%d %s.%s %s[%s] %s\n",
		cTime, cDecimal[:6], "INFO ", filepath.Base(cFile), cLine, cPkg, cFunc,
		cPrefix, cContextPair, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)
}

func TestSyslogHeader(t *testing.T) {
	formatter := text.New(text.NewConfig().WithHeader(text.SyslogHeader))
	expect := fmt.Sprintf("%s:%d %s.%s %s[%s] %s\n",
		filepath.Base(cFile), cLine, cPkg, cFunc, cPrefix, cContextPair, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)
}

func TestCustomHeader(t *testing.T) {
	header := "{{time:time.ns}} {{level:char}} {{file}}:{{line%05d}} {{pkg:1}}.{{func}} " +
		"{{prefix}}[{{context:list}}] {{msg%20s}}\n"
	formatter := text.New(text.NewConfig().WithHeader(header))
	expect := fmt.Sprintf("%s.%s %s %s:%05d %s.%s %s[%s] %20s\n",
		cTime, cDecimal, "I", cFile, cLine, path.Base(cPkg), cFunc,
		cPrefix, cContextList, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)
}

func TestBizarreHeader(t *testing.T) {
	formatter := text.New(text.NewConfig())
	header := "xx{{unknown}} {static} {{unknown}} {{level : char}} {{pkg|1}} " +
		"[{{context:dot}}] {{ msg %20s }}yy"
	formatter.SetHeader(header)
	expect := fmt.Sprintf("xx {static}  %s  [%s] %20syy",
		"I", cContextPair, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)
}

func TestColor(t *testing.T) {
	formatter := text.New(text.NewConfig().WithHeader("{{msg}}").WithEnableColor(true))
	expect := fmt.Sprintf("\033[%dm%s\033[0m", text.Magenta, cMsg)
	testFormat(t, formatter, &gTmplRecord, expect)

	record := cloneRecord()
	record.Level = gxlog.LevelWarn
	record.Aux.Marked = false
	formatter.MapColors(map[gxlog.Level]text.ColorID{
		gxlog.LevelWarn: text.Blue,
	})
	expect = fmt.Sprintf("\033[%dm%s\033[0m", text.Blue, cMsg)
	testFormat(t, formatter, record, expect)

	record.Level = gxlog.LevelError
	formatter.SetColor(gxlog.LevelError, text.Yellow)
	expect = fmt.Sprintf("\033[%dm%s\033[0m", text.Yellow, cMsg)
	testFormat(t, formatter, record, expect)
}

func testFormat(t *testing.T, formatter gxlog.Formatter, record *gxlog.Record, expect string) {
	output := string(formatter.Format(record))
	if output != expect {
		t.Errorf("testFormat:\noutput: %q\nexpect: %q", output, expect)
	}
}

func cloneRecord() *gxlog.Record {
	clone := gTmplRecord
	clone.Aux.Contexts = make([]gxlog.Context, len(gTmplRecord.Aux.Contexts))
	copy(clone.Aux.Contexts, gTmplRecord.Aux.Contexts)
	return &clone
}
