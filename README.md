# gxlog #

Gxlog is short for **G**o e**X**tensible **LOG**ger. It is concise, functional,
flexible and extensible. Easy-to-use is also an important design goal. It only
depends on the standard library.

Besides the basic functionality of logging, gxlog also provides many advanced
features, such as static context, dynamic context, log limitation and so on.
With the interface `Formatter` and `Writer`, gxlog can be extended to support
any log format or backend. In addition, with the design of **Slot**, logging,
eventing or hooking can be integrated into gxlog.

## Contents ##

- [Architecture](#architecture)
- [Features Preview](#features-preview)
- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Basic](#basic)
  - [Auxiliary](#auxiliary)
  - [Slots](#slots)
  - [Settings](#settings)
  - [New Logger](#new-logger)
  - [Formatters](#formatters)
  - [Writers](#writers)

## Architecture ##

```
+-------------------------------------------------+
|                     logger                      |
|            level [filter] [limiter]             |
| record                                          |
|   | +-------+---------------------------------+ |
|   |-| Slot0 | formatter writer level [filter] | |
|   | +-------+---------------------------------+ |
|   |-|  ...  |              ...                | |
|   | +-------+---------------------------------+ |
|   \-| Slot7 | formatter writer level [filter] | |
|     +-------+---------------------------------+ |
+-------------------------------------------------+
```

A logger contains EIGHT slots. Each slot contains a formatter and a writer.
The logger has its own level and filter while each slot has its independent
level and filter. When a log record is emitted, the logger calls the formatter
and writer of each slot in the order from Slot0 to Slot7 to format and write
the log.

## Features Preview ##

- **logger**
  - level
  - filter
  - prefix
  - context
  - mark
  - limitation
  - helper methods
  - auto backtracking
  - **slots**
    - manipulation
    - level
    - filter
  - **formatter**
    - formatter function wrapper
    - null formatter
    - **text formatter**
      - custom property and format of fields
      - colorization
      - custom color mapping
    - **json formatter**
      - custom property of fields
      - custom omission of fields
      - custom omission of empty fields
  - **writer**
    - writer function wrapper
    - multi-writer wrapper
    - io.Writer wrapper
    - asynchronous wrapper
    - null writer
    - **file writer**
      - custom file naming
      - file splitting
      - file deletion checking
      - new directory each day
      - gzip compression
      - AES encryption
      - error handler
    - **syslog writer**
      - custom mapping from level to severity
      - error handler
    - **tcp socket writer**
    - **unix domain socket writer**

## Getting Started ##

### Installing ###

To install gxlog, run `go get`:

``` shell
$ go get github.com/gxlog/gxlog/...
```

### Basic ###

The default Logger has the default Formatter (a text formatter) and a writer
wrapper of os.Stderr linked at Slot0. The rest slots are free.

It is **RECOMMENDED** that all packages use the default Logger, such the main
package can control which, how and where to output logs by setting filters,
formatters and writers of the default Logger.

Supported levels are Trace, Debug, Info, Warn, Error and Fatal. Timing and
error helper methods are provided.

All methods of a Logger are concurrency safe.

``` go
package main

import (
    "time"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/iface"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
    // gxlog.Formatter returns the default Formatter in Slot0.
    // Coloring is only supported on systems that ANSI escape sequences
    // are supported.
    gxlog.Formatter().EnableColor()

    testLevel()
    // testPanic()
    testTime()
    testLog()
    testLogError()
    testLogErrorf()
}

func testLevel() {
    log.Trace("test Trace")
    log.Tracef("%s", "test Tracef")
    log.Debug("test Debug")
    log.Debugf("%s", "test Debugf")
    log.Info("test Info")
    log.Infof("%s", "test Infof")
    log.Warn("test Warn")
    log.Warnf("%s", "test Warnf")
    log.Error("test Error")
    log.Errorf("%s", "test Errorf")
    // Fatal and Fatalf will output the stack of current goroutine by default.
    log.Fatal("test Fatal")
    log.Fatalf("%s", "test Fatalf")
}

func testPanic() {
    // The default Level of Panic or Panicf is Fatal.
    // It will always panic when Panic or Panicf is called no matter at which
    // level the Logger is.
    log.Panic("test Panic")
    log.Panicf("%s", "test Panicf")
}

func testTime() {
    // Time or Timef returns a function. When the function is called, it outputs
    // the log as well as the time cost since the call of Time or Timef.
    // The default Level of Time or Timef is Trace.
    done := log.Time("test Time")
    time.Sleep(200 * time.Millisecond)
    done()
    // Time or Timef works well with defer.
    // Notice the last empty pair of parentheses.
    defer log.Timef("%s", "test Timef")()
    time.Sleep(400 * time.Millisecond)
}

func testLog() {
    // The calldepth can be specified in Log and Logf. That is useful when you
    // are customizing your own log wrapper function.
    log.Log(0, iface.Info, "test Log")
    log.Logf(1, iface.Warn, "%s: %d", "test Logf", 1)
    log.Logf(-1, iface.Warn, "%s: %d", "test Logf", -1)
}

func testLogError() error {
    // LogError outputs a log and call errors.New to generate an error.
    return log.LogError(iface.Error, "an error")
}

func testLogErrorf() error {
    // LogErrorf outputs a log and call fmt.Errorf to generate an error.
    return log.LogErrorf(iface.Error, "%s", "another error")
}
```

### Auxiliary ###

The methods with the prefix `With` of Logger will attach auxiliary information
to log records or limit log output. They can be chained together in any number.
Calls to `WithContext` will concatenate the context key-value pairs while calls
to the others will overwrite corresponding settings.

In fact, each call to any of them will return a new instance of Logger which is
a shallow copy of the calling Logger. The prefix, contexts, mark, count limiter
and time limiter of Logger are copied before they are modified. Such, instances
of Logger have the **lexical scope**.

``` go
package main

import (
    "math/rand"
    "time"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/logger"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func init() {
    rand.Seed(time.Now().UnixNano())
}

func main() {
    // gxlog.Formatter returns the default Formatter in Slot0.
    // Coloring is only supported on systems that ANSI escape sequences
    // are supported.
    gxlog.Formatter().EnableColor()

    testAuxiliary()
    testDynamicContext()
    testLimitation()
}

func testAuxiliary() {
    // Logs with mark will be colorized with Magenta by default.
    // The prefix or mark allow you to highlight some logs temporarily
    // while you are debugging.
    log.WithPrefix("**** ").WithMark(true).WithContext("k1", "v1", "k2", "v2").
        Info("prefix, mark and contexts")
    // The original log instance is not altered.
    log.Info("no prefix, mark or contexts")

    // This demonstrates the lexical scope of a log instance:
    func() {
        log := log.WithContext("k3", "v3")
        log.Info("outer enter")
        func() {
            log := log.WithContext("k4", "v4")
            log.Info("inner")
        }()
        log.Info("outer leave")
    }()
}

func testDynamicContext() {
    // All the key-value pairs of dynamic contexts will be appended to the end
    // of static contexts.
    // Dynamic contexts are very useful when you want to print the current value
    // of some variables all the time.
    // ATTENTION: You SHOULD be very careful to concurrency safety or deadlocks
    // with dynamic contexts.
    n := 0
    fn := logger.Dynamic(func(interface{}) interface{} {
        // Do NOT call any method of the Logger in the function,
        // or it may deadlock.
        n++
        return n
    })
    clog := log.WithContext("static", n, "dynamic", fn)
    clog.Info("dynamic one")
    clog.Info("dynamic two")
}

func testLimitation() {
    // THINK TWICE before you decide to limit the output of logs by count or
    // by time, you may miss logs which you need.
    // Only 2 logs will be output per 3 logs.
    for i := 1; i <= 6; i++ {
        log.WithCountLimit(3, 2).Infof("count limited: %d", i)
    }
    // the more efficient way
    llog := log.WithCountLimit(3, 2)
    for i := 7; i <= 12; i++ {
        llog.Infof("efficient count limited: %d", i)
    }
    // NOTICE: The space complexity is O(n), while n is the 2nd argument of
    // WithTimeLimit. Try to specify reasonable duration and limit.
    // At most 3 logs will be output during any interval of 1 second.
    for i := 1; i <= 10; i++ {
        log.WithTimeLimit(time.Second, 3).Infof("time limited: %d", i)
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    }
}
```

### Slots ###

A Logger has EIGHT slots. The Formatter and Writer in each slot will be called in
the order from `Slot0` to `Slot7` when a log is emitted. Custom formatters or
writers can act as event triggers or hooks. Each slot has its independent Level
and Filter.

You can do linking or unlinking at any slot, copy or move a slot to another or
swap any two slots. The Formatter, Writer, Level and Filter of a slot can be set
independently at any time.

``` go
package main

import (
    "fmt"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/formatter"
    "github.com/gxlog/gxlog/formatter/json"
    "github.com/gxlog/gxlog/iface"
    "github.com/gxlog/gxlog/logger"
    "github.com/gxlog/gxlog/writer"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
    // gxlog.Formatter returns the default Formatter in Slot0.
    // Coloring is only supported on systems that ANSI escape sequences
    // are supported.
    gxlog.Formatter().EnableColor()

    testSlots()
    testSlotsLevel()
}

func testSlots() {
    log.Info("this will be printed once")

    // Copy the Formatter, Writer, Filter and Level of Slot0 to Slot1.
    log.CopySlot(logger.Slot1, logger.Slot0)
    log.Info("this will be printed twice")

    log.SetSlotFormatter(logger.Slot1, json.New(json.Config{}))
    log.Info("this will be printed in text format and json format")

    log.SwapSlot(logger.Slot0, logger.Slot1)
    log.Info("json first and then text")

    // Copy the Formatter, Writer, Filter and Level of Slot1 to Slot0 and then
    // set the Formatter to formatter.Null(), Writer to writer.Null(), Filter to
    // nil and Level to Off of Slot1.
    log.MoveSlot(logger.Slot0, logger.Slot1)
}

func testSlotsLevel() {
    log.SetSlotLevel(logger.Slot0, iface.Warn)
    log.Info("this will not be printed")
    log.Warn("this will be printed")

    log.SetSlotLevel(logger.Slot0, iface.Trace)
    // A Formatter or Writer can act as a hook.
    // ATTENTION: Do NOT call any method of the Logger in a Formatter, Writer
    // or Filter, or it may deadlock.
    hook := formatter.Func(func(record *iface.Record) []byte {
        // log.Info("deadlock")
        fmt.Println("hooks:", record.Msg)
        return nil
    })
    filter := func(record *iface.Record) bool {
        return record.Aux.Marked
    }
    // Link at Slot0 will overwrite the current link at Slot0.
    // Use formatter.Null() instead of a nil Formatter and writer.Null()
    // instead of a nil Writer, or it will panic.
    // If the Level of a log is NOT lower than Warn and it is marked, the hook
    // will be called.
    log.Link(logger.Slot0, hook, writer.Null(), iface.Warn, filter)
    log.WithMark(true).Info("marked, but info")
    log.Error("error, but not marked")
    log.WithMark(true).Warn("warn and marked")
}
```

### Settings ###

The logger has a bundle of methods to get and set its levels, flags and filter.
They are all concurrency safe and you can alter the settings at any time.

The function `And`, `Or` and `Not` allow you to combine filters into a more
complex one.

``` go
package main

import (
    "strings"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/iface"
    "github.com/gxlog/gxlog/logger"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
    // gxlog.Formatter returns the default Formatter in Slot0.
    // Coloring is only supported on systems that ANSI escape sequences
    // are supported.
    gxlog.Formatter().EnableColor()

    testConfig()
    testFilterLogic()
}

func testConfig() {
    log.Infof("config: %#v", log.Config())

    log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).
        Fatal("fatal before updating the config")
    log.UpdateConfig(func(config logger.Config) logger.Config {
        // Do NOT call any method of the Logger in the function,
        // or it may deadlock.
        // Disable Prefix, StaticContext and Mark, then their value will always
        // be the zero value of their type.
        config.Disabled |= (logger.Prefix | logger.StaticContext | logger.Mark)
        // Disable auto backtracking
        config.TrackLevel = iface.Off
        return config
    })
    log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).
        Fatal("fatal after updating the config")
}

func testFilterLogic() {
    log.SetFilter(logger.Or(important, logger.And(useful, interesting)))
    log.Error("error") // this will be output
    log.Warn("warn")
    log.Trace("trace, funny")
    log.Info("info, funny") // this will be output
}

func important(record *iface.Record) bool {
    return record.Level >= iface.Error
}

func useful(record *iface.Record) bool {
    return record.Level >= iface.Info
}

func interesting(record *iface.Record) bool {
    return strings.Contains(record.Msg, "funny")
}
```

### New Logger ###

If you really need a new Logger rather than the default one, you can create it.

``` go
package main

import (
    "fmt"
    "os"

    "github.com/gxlog/gxlog/formatter/json"
    "github.com/gxlog/gxlog/formatter/text"
    "github.com/gxlog/gxlog/iface"
    "github.com/gxlog/gxlog/logger"
    "github.com/gxlog/gxlog/writer"
    "github.com/gxlog/gxlog/writer/file"
)

func main() {
    log := logger.New(logger.Config{
        Disabled:   logger.DynamicContext | logger.LimitByTime,
        TrackLevel: iface.Off,
    })

    fileWriter, err := file.Open(file.Config{
        Path: "/tmp/gxlog",
        Base: "test",
        // ErrorHandler will be called when an error occurs.
        ErrorHandler: writer.Report,
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer fileWriter.Close()

    // Logs will be formatted to text format and output to os.Stderr, then
    // formatted to json format and output to log files in /tmp/gxlog.
    log.Link(logger.Slot0, text.New(text.Config{
        EnableColor: true,
    }), writer.Wrap(os.Stderr))
    log.Link(logger.Slot1, json.New(json.Config{}), fileWriter)

    log.Info("test")
}
```

### Formatters ###

Gxlog provides a function wrapper, a text formatter and a json formatter. The
text formatter has support for colorization and its header is highly customizable.
The json formatter can be customized to omit specified fields of a log.

All methods of a text formatter or a json formatter are concurrency safe and you
can alter the config of a formatter at any time.

``` go
package main

import (
    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/formatter"
    "github.com/gxlog/gxlog/formatter/json"
    "github.com/gxlog/gxlog/formatter/text"
    "github.com/gxlog/gxlog/iface"
    "github.com/gxlog/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
    testCustomFormatter()
    testTextFormatter()
    testJSONFormatter()
}

func testCustomFormatter() {
    fn := formatter.Func(func(record *iface.Record) []byte {
        return []byte(record.Msg + "\n")
    })
    log.SetSlotFormatter(logger.Slot0, fn)
    log.Info("a simple formatter that just returns the msg of a record")
}

func testTextFormatter() {
    // By default, Trace, Debug and Info map to Green, Warn maps to Yellow,
    // Error and Fatal map to Red, marked logs map to Magenta.
    textFmt := text.New(text.Config{
        // Coloring is only supported on systems that ANSI escape sequences
        // are supported.
        EnableColor: true,
        Header:      text.CompactHeader,
    })
    log.SetSlotFormatter(logger.Slot0, textFmt)
    log.Trace("green")
    log.Warn("yellow")
    log.Error("red")
    log.WithMark(true).Error("magenta")

    // update settings
    textFmt.SetHeader(text.FullHeader)
    textFmt.SetColor(iface.Trace, text.Blue)
    textFmt.MapColors(map[iface.Level]text.Color{
        iface.Warn:  text.Red,
        iface.Error: text.Magenta,
    })
    textFmt.SetMarkedColor(text.White)
    log.Trace("blue")
    log.Warn("red")
    log.Error("magenta")
    log.WithMark(true).Error("white")

    // For details of all supported fields, see the comment of text.Config.
    header := "{{time:time}} {{level:char}} {{file:2%q}}:{{line:%05d}} {{msg:%20s}}\n"
    textFmt.SetHeader(header)
    textFmt.DisableColor()
    log.Trace("default color")
}

func testJSONFormatter() {
    jsonFmt := json.New(json.Config{
        // Only the last segment of the File field will be formatted.
        FileSegs: 1,
    })
    log.SetSlotFormatter(logger.Slot0, jsonFmt)
    log.Trace("json")

    // update settings
    jsonFmt.UpdateConfig(func(config json.Config) json.Config {
        // Do NOT call any method of the Formatter or the Logger in the function,
        // or it may deadlock.
        config.OmitEmpty = json.Aux
        config.Omit = json.Pkg | json.Func
        return config
    })
    log.Trace("json updated")
    log.WithContext("ah", "ha").Trace("json with context")
}
```

### Writers ###

Gxlog provides several writer wrappers to gxlog.Writer and several writers.

Gxlog provides several writer wrappers, including function wrapper, io.Writer
wrapper, multi-writer wrapper and asynchronous wrapper. The multi-writer wrapper
can combine several writers into one if they are writing logs with the same
format. The asynchronous wrapper can make a writer switch into asynchronous mode.

Gxlog also provides several writers including tcp socket writer, unix domain
socket writer, file writer and syslog writer. Both the tcp socket writer and the
unix domain socket writer aim at log watching and make 'netcat' a alternative to
the 'tail'. Tailing a file is somewhat inconvenient because a new log file will
be created when a log file reaches its max size. For log transmission, use the
syslog writer instead. You can register an error handler for file writer and
syslog writer.

For performance and security, use the unix socket domain writer instead of the
tcp socket writer as long as the system has support for unix domain socket.
Otherwise, when using a tcp socket writer, bind the address to localhost only.

For performance and security, connect to the local syslog server and configure
the local syslog server for log transmission if it is possible.

All methods of a writer are concurrency safe and you can alter the config of a
writer at any time.

``` go
package main

import (
    "compress/flate"
    "fmt"
    "os"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/formatter/text"
    "github.com/gxlog/gxlog/iface"
    "github.com/gxlog/gxlog/logger"
    "github.com/gxlog/gxlog/writer"
    "github.com/gxlog/gxlog/writer/file"
    "github.com/gxlog/gxlog/writer/socket/tcp"
    "github.com/gxlog/gxlog/writer/socket/unix"
    "github.com/gxlog/gxlog/writer/syslog"
)

// gxlog.Logger returns the default Logger.
var log = gxlog.Logger()

func main() {
    // gxlog.Formatter returns the default Formatter in Slot0.
    // Coloring is only supported on systems that ANSI escape sequences
    // are supported.
    gxlog.Formatter().EnableColor()

    testWrappers()
    testSocketWriters()

    gxlog.Formatter().DisableColor()

    testFileWriter()
    testSyslogWriter()
}

func testWrappers() {
    // custom writer function
    fn := writer.Func(func(bs []byte, _ *iface.Record) {
        os.Stderr.Write(bs)
    })
    log.SetSlotWriter(logger.Slot0, fn)
    log.Info("a simple writer that just writes to os.Stderr")

    // another equivalent way
    log.SetSlotWriter(logger.Slot0, writer.Wrap(os.Stderr))
    log.Info("writer wrapper of os.Stderr")

    // multi-writer
    multi := writer.Multi(writer.Wrap(os.Stdout), writer.Wrap(os.Stderr))
    log.SetSlotWriter(logger.Slot0, multi)
    log.Info("multi-writer: this will be printed twice")

    // Asynchronous writer wrapper uses a internal channel to buffer logs.
    // When the channel is full, the Write method of the wrapper blocks.
    // ATTENTION: Some logs may NOT be output in asynchronous mode if os.Exit
    // is called, panicking without recovery and so on.
    async := writer.NewAsync(writer.Wrap(os.Stderr), 1024)
    // Close waits until all logs in the channel have been output.
    // It does NOT close the underlying writer.
    // To ignore all logs that have not been output, use Abort instead.
    defer async.Close()

    log.SetSlotWriter(logger.Slot0, async)
    log.Info("asynchronous writer wrapper")
}

func testSocketWriters() {
    // tcp socket writer
    // The default address is "localhost:9999".
    tcpWriter, err := tcp.Open(tcp.Config{})
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer tcpWriter.Close()

    log.SetSlotWriter(logger.Slot0, tcpWriter)

    // Use `netcat localhost 9999' to watch logs.
    // for i := 0; i < 1024; i++ {
    //  log.Info(i)
    //  time.Sleep(time.Second)
    // }

    // unix domain socket writer
    // Shell expansion is NOT supported. Thus, ~, $var and so on will NOT
    // be expanded.
    // The default pathname is "/tmp/gxlog/<pid>".
    unixWriter, err := unix.Open(unix.Config{
        Pathname: "/tmp/gxlog/unixdomain",
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer unixWriter.Close()

    log.SetSlotWriter(logger.Slot0, unixWriter)

    // Use "netcat -U /tmp/gxlog/unixdomain" to watch logs.
    // for i := 0; i < 1024; i++ {
    //  log.Info(i)
    //  time.Sleep(time.Second)
    // }
}

func testFileWriter() {
    gxlog.Formatter().SetHeader(text.CompactHeader)

    // Shell expansion is NOT supported. Thus, ~, $var and so on will NOT
    // be expanded.
    wt, err := file.Open(file.Config{
        Path:      "/tmp/gxlog",
        Base:      "test",
        DateStyle: file.DateUnderscore,
        TimeStyle: file.TimeDot,
        // ErrorHandler will be called when an error occurs.
        ErrorHandler: writer.Report,
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer wt.Close()

    log.SetSlotWriter(logger.Slot0, wt)
    log.Info("this will be output to a file")

    wt.UpdateConfig(func(config file.Config) file.Config {
        // Do NOT call any method of the Writer or the Logger in the function,
        // or it may deadlock.
        config.Ext = ".bin"
        // enable gzip compression
        config.GzipLevel = flate.DefaultCompression
        // enable AES encryption.
        // The key MUST be hexadecimal encoded without the prefix 0X or 0x.
        config.AESKey = "70856575b161fbcca8fc12e1f70fc1c8"
        return config
    })
    log.Info("gzipped and encrypted")
}

func testSyslogWriter() {
    gxlog.Formatter().SetHeader(text.SyslogHeader)

    // Leave the Network field of the config to be empty, it will connect to the
    // local syslog server with unix domain socket.
    wt, err := syslog.Open(syslog.Config{
        Tag:      "gxlog",
        Facility: syslog.FacUser,
        // ErrorHandler will be called when an error occurs.
        ErrorHandler: writer.Report,
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer wt.Close()

    log.SetSlotWriter(logger.Slot0, wt)
    // NOTICE: The standard syslog package will get the timestamp itself which is
    // a tiny bit later than Record.Time.
    log.Info("this will be output to syslog")

    // update level mapping
    // The severity of a level is left to be unchanged if it is not in the map.
    wt.MapSeverity(map[iface.Level]syslog.Severity{
        iface.Info: syslog.SevErr,
    })
    log.Info("this will be severity err")
}
```
