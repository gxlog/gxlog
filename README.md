# gxlog #

Gxlog is short for **G**o e**X**tensible **LOG**ger. It is concise, functional,
flexible and extensible. Easy-to-use is also an important design goal. Besides
the basic functionality of logging, gxlog also provides many advanced features,
such as context, dynamic context, log limitation and so on. With the interface
**Formatter** and **Writer**, gxlog can be extended to support any log formats
or any log backends. In addition, with the design of **Slots**, logging, events
or hooks can be integrated into gxlog.

## Table of Contents ##

- [Architecture](#architecture)
- [Features Preview](#features-preview)
- [Getting Started](#getting-started)
  - [Basic](#basic)
  - [Auxiliary](#auxiliary)
  - [Slots](#slots)
  - [Settings](#settings)
  - [New Logger](#new-logger)
  - [Formatters](#formatters)

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

## Features Preview ##

- **logger**
  - level
  - filter
  - filter logic
  - prefix
  - context
  - dynamic context
  - mark
  - limitation
  - timekeeping helper
  - error helper
  - auto backtrack
  - **slots**
    - management
    - level
    - filter
  - **formatter**
    - formatter function wrapper
    - **text formatter**
      - custom header
      - custom property and format of fields
      - colorization
      - custom color mapping
    - **json formatter**
      - custom property of fields
      - custom omission of fields
      - custom omission of empty fields
  - **writer**
    - writer function wrapper
    - io.Writer wrapper
    - asynchronous wrapper
    - **file writer**
      - custom file max size
      - custom file naming style
      - file deletion check
      - new directory each day
      - gzip compression
      - AES encryption
      - error reporting
    - **syslog writer**
      - custom mapping from level to severity
      - error reporting
    - **tcp socket writer**
    - **unix domain socket writer**

## Getting Started ##

### Basic ###

The default Logger has the default Formatter (a text formatter) and a writer
wrapper of os.Stderr linked in Slot0. The rest slots are free.

It is **RECOMMENDED** that all packages use the default logger, such the main
package can control which, how and where to output logs by setting filters,
formatters and writers of the default logger.

Supported levels are TRACE, DEBUG, INFO, WARN, ERROR and FATAL.

All methods of a Logger are concurrency safe.

``` go
package main

import (
    "time"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func main() {
    // Only supported on systems that ANSI escape sequences are supported.
    defaults.Formatter().EnableColor()

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

    // The default level of Panic or Panicf is fatal.
    // It will always panic no matter at which level the logger is.
    // log.Panic("test Panic")
    // log.Panicf("%s", "test Panicf")

    // Time and Timef will return a function. When the function is called,
    //   it will output the log as well as the time cost since the call of
    //   Time or Timef.
    // The default level of Time and Timef is trace.
    done := log.Time("test Time")
    time.Sleep(200 * time.Millisecond)
    done()
    // Notice the last empty pair of parentheses.
    defer log.Timef("%s", "test Timef")()
    time.Sleep(400 * time.Millisecond)

    // The calldepth can be specified in Log and Logf. That is useful when
    //   you want to customize your own log helper functions.
    log.Log(0, gxlog.Info, "test Log")
    log.Logf(1, gxlog.Warn, "%s: %d", "test Logf", 1)
    log.Logf(-1, gxlog.Warn, "%s: %d", "test Logf", -1)

    test1()
    test2()
}

func test1() error {
    // LogError will output log and call errors.New to generate an error
    return log.LogError(gxlog.Error, "an error")
}

func test2() error {
    // LogErrorf will output log and call fmt.Errorf to generate an error
    return log.LogErrorf(gxlog.Error, "%s", "another error")
}
```

### Auxiliary ###

The methods with the prefix `With` of Logger will attach auxiliary information
to log records or limit log output. They can be chained together in any number.
Calls to `WithContext` will concatenate the context key-value pairs while calls
to the others will overwrite corresponding settings.

In fact, each call to any of them will return a new instance of Logger which is
a shallow copy of the caller. The prefix, contexts, mark, count limiter and time
limiter of Logger are copied before they are modified. Such, instances of Logger
have the **lexical scope**.

``` go
package main

import (
    "math/rand"
    "time"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func init() {
    rand.Seed(time.Now().UnixNano())
}

func main() {
    // Only supported on systems that ANSI escape sequences are supported.
    defaults.Formatter().EnableColor()

    // logs with mark will be colorized with magenta by default
    // prefix and mark allow you to highlight some logs while you are debugging
    log.WithPrefix("**** ").WithMark(true).WithContext("k1", "v1", "k2", "v2").
        Info("prefix, mark and contexts")
    // the instance of log is left to be unchanged
    log.Info("no prefix, mark or contexts")

    // demonstrates the lexical scope
    func() {
        log := log.WithContext("k3", "v3")
        log.Info("outer enter")
        func() {
            log := log.WithContext("k4", "v4")
            log.Info("inner")
        }()
        log.Info("outer leave")
    }()

    // all the key-value pairs of dynamic contexts will be concatenated to the
    //   end of static contexts
    // dynamic contexts are very useful when you want to print some/all fields
    //   of a struct value all the time.
    // ATTENTION: you should be very careful to concurrency safety or deadlocks
    //   with dynamic contexts.
    n := 0
    fn := gxlog.Dynamic(func(interface{}) interface{} {
        n++
        return n
    })
    clog := log.WithContext("static", n, "dynamic", fn)
    clog.Info("dynamic one")
    clog.Info("dynamic two")

    // THINK TWICE before you limit logs output by count or by time,
    //   you may miss logs which you need.
    // only 2 logs will be output per 3 logs
    for i := 1; i <= 6; i++ {
        log.WithCountLimit(3, 2).Infof("count limited: %d", i)
    }
    // the more efficient way
    llog := log.WithCountLimit(3, 2)
    for i := 7; i <= 12; i++ {
        llog.Infof("efficient count limited: %d", i)
    }
    // NOTICE: The space complexity is O(n), while n is the 2nd argument of
    //   WithTimeLimit. Try to specify reasonable duration and limit.
    // at most 3 logs will be output during any interval of 1 second
    for i := 1; i <= 10; i++ {
        log.WithTimeLimit(time.Second, 3).Infof("time limited: %d", i)
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    }
}
```

### Slots ###

The logger has 8 slots, from `Slot0` to `Slot7`. The formatter and writer in each
slot will be called in order from `Slot0` to `Slot7`. Custom formatters or writers
can act as event triggers or hooks. Each slot has independent level and filter.

``` go
package main

import (
    "fmt"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/defaults"
    "github.com/gxlog/gxlog/formatter"
    "github.com/gxlog/gxlog/formatter/json"
)

var log = defaults.Logger()

func main() {
    // Only supported on systems that ANSI escape sequences are supported.
    defaults.Formatter().EnableColor()

    log.Info("this will print once")

    // copy Slot0 with the default formatter and wrapper of os.Stderr to Slot1
    log.CopySlot(gxlog.Slot1, gxlog.Slot0)
    log.Info("this will print twice")

    log.SetSlotFormatter(gxlog.Slot1, json.New(json.NewConfig()))
    log.Info("this will print in text format and json format")

    log.SwapSlot(gxlog.Slot0, gxlog.Slot1)
    log.Info("json first and then text")

    // set the formatter, writer and filter of Slot0 to nil and
    //   set the level of Slot0 to off
    log.Unlink(gxlog.Slot0)

    log.SetSlotLevel(gxlog.Slot1, gxlog.Warn)
    log.Info("this will not print")
    log.Warn("this will print")

    log.SetSlotLevel(gxlog.Slot1, gxlog.Trace)
    // ATTENTION: DO NOT call methods of logger in formatter, writer or filter
    //   in the current goroutine, or it will deadlock.
    hook := formatter.Func(func(record *gxlog.Record) []byte {
        // log.Info("deadlock")
        fmt.Println("hooks:", record.Msg)
        return nil
    })
    filter := func(record *gxlog.Record) bool {
        return record.Aux.Marked
    }
    // link at Slot0 will overwrite the current link at Slot0 if any
    // If the log level is not lower than WARN and the log is marked, the hook
    //   will be called.
    log.Link(gxlog.Slot0, hook, nil, gxlog.Warn, filter)
    log.WithMark(true).Info("marked, but info")
    log.Error("error, but not marked")
    log.WithMark(true).Warn("warn and marked")
}
```

### Settings ###

The logger has a bundle of methods to get or set different levels, flags or the
filter. They are all concurrency safe and you can alter the config of the logger
at any time.

``` go
package main

import (
    "strings"

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/defaults"
)

var log = defaults.Logger()

func main() {
    // Only supported on systems that ANSI escape sequences are supported.
    defaults.Formatter().EnableColor()

    log.Infof("config: %#v", log.Config())

    log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).Fatal("fatal before update")
    log.UpdateConfig(func(config gxlog.Config) gxlog.Config {
        // Do NOT call methods of the Logger, or it will deadlock.
        // disable prefix, contexts and mark
        // these attributes of records will always be the zero value of their type
        config.Flags &^= (gxlog.Prefix | gxlog.Contexts | gxlog.Mark)
        // disable the auto backtracking
        config.TrackLevel = gxlog.Off
        return config
    })
    log.WithPrefix("**** ").WithContext("k1", "v1").WithMark(true).Fatal("fatal after update")

    // demonstrates the filter logic
    log.SetFilter(gxlog.Or(important, gxlog.And(useful, interesting)))
    log.Error("error") // this will print
    log.Warn("warn")
    log.Trace("trace, funny")
    log.Info("info, funny") // this will print
}

func important(record *gxlog.Record) bool {
    return record.Level >= gxlog.Error
}

func useful(record *gxlog.Record) bool {
    return record.Level >= gxlog.Info
}

func interesting(record *gxlog.Record) bool {
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

    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/formatter/json"
    "github.com/gxlog/gxlog/formatter/text"
    "github.com/gxlog/gxlog/writer"
    "github.com/gxlog/gxlog/writer/file"
)

func main() {
    // create a new Logger with default config
    log := gxlog.New(gxlog.NewConfig())

    // create a new Config and customize it
    // config := gxlog.NewConfig().
    //  WithDisabled(gxlog.DynamicContexts | gxlog.Limit).
    //  WithTrackLevel(gxlog.Off)

    // another equivalent way
    // config := gxlog.NewConfig()
    // config.Flags &^= gxlog.DynamicContexts | gxlog.Limit
    // config.TrackLevel = gxlog.Off

    // create a new Logger with custom config
    // gxlog.New(config)

    // create a file writer, logs output to /tmp/gxlog
    fileWriter, err := file.Open(file.NewConfig("/tmp/gxlog", "base"))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer fileWriter.Close()

    log.Link(gxlog.Slot0, text.New(text.NewConfig().WithEnableColor(true)),
        writer.Wrap(os.Stderr))
    log.Link(gxlog.Slot1, json.New(json.NewConfig()), fileWriter)

    log.Info("test")
}
```

### Formatters ###

Gxlog provides a function wrapper to gxlog.Formatter, a text formatter and a json
formatter. You can create a formatter with custom config or just leave it to the
default config and update the settings of the formatter later.

All methods of a text formatter or a json formatter are concurrency safe and you
can alter the config of a formatter at any time.

``` go
package main

import (
    "github.com/gxlog/gxlog"
    "github.com/gxlog/gxlog/defaults"
    "github.com/gxlog/gxlog/formatter"
    "github.com/gxlog/gxlog/formatter/json"
    "github.com/gxlog/gxlog/formatter/text"
)

var log = defaults.Logger()

func main() {
    // custom formatter function
    fn := formatter.Func(func(record *gxlog.Record) []byte {
        return append([]byte(record.Msg), '\n')
    })
    log.SetSlotFormatter(gxlog.Slot0, fn)
    log.Info("a simple formatter that just returns record.Msg")

    // text formatter
    textFmt := text.New(text.NewConfig().
        WithEnableColor(true).
        WithHeader(text.CompactHeader))
    log.SetSlotFormatter(gxlog.Slot0, textFmt)
    log.Trace("green")
    log.Warn("yellow")
    log.Error("red")
    log.WithMark(true).Error("magenta")

    // update settings of the text formatter
    textFmt.SetHeader(text.DefaultHeader)
    textFmt.SetColor(gxlog.Trace, text.Blue)
    textFmt.MapColors(map[gxlog.Level]text.ColorID{
        gxlog.Warn:  text.Red,
        gxlog.Error: text.Magenta,
    })
    textFmt.SetMarkedColor(text.White)
    log.Trace("blue")
    log.Warn("red")
    log.Error("magenta")
    log.WithMark(true).Error("white")

    // custom header of text formatter
    textFmt.SetHeader("{{time:time}} {{level:char}} {{file:2%q}}:{{line:%05d}} {{msg:%20s}}\n")
    textFmt.DisableColor()
    log.Trace("default color")

    // json formatter
    jsonFmt := json.New(json.NewConfig().WithFileSegs(1))
    log.SetSlotFormatter(gxlog.Slot0, jsonFmt)
    log.Trace("json")

    // update settings of the json formatter
    jsonFmt.UpdateConfig(func(config json.Config) json.Config {
        config.OmitEmpty = json.Aux
        config.Omit = json.Pkg | json.Func
        return config
    })
    log.Trace("json updated")
    log.WithContext("ah", "ha").Trace("json with contexts")
}
```

