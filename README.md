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
    - null formatter
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
    - null writer
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

The default logger has text formatter and writer wrapper to os.Stderr linked in
Slot0. The rest slots are free. Supported levels are trace, debug, info, warn,
error and fatal.

It is RECOMMENDED that all packages use the default logger, such the main package
can control which, how and where to output logs by setting filters, formatters
and writers of the default logger.

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
    //   it will output the log as well as the time cost.
    // The default level of Time and Timef is trace.
    done := log.Time("test Time")
    time.Sleep(200 * time.Millisecond)
    done()
    // Notice the last empty pair of parentheses.
    defer log.Timef("%s", "test Timef")()
    time.Sleep(400 * time.Millisecond)

    // The calldepth can be specified in Log and Logf. That is useful when
    //   you want to customize your own log helper functions.
    log.Log(0, gxlog.LevelInfo, "test Log")
    log.Logf(-1, gxlog.LevelWarn, "%s", "test Logf")

    test1()
    test2()
}

func test1() error {
    // LogError will output log and call errors.New to generate an error
    return log.LogError(gxlog.LevelError, "an error")
}

func test2() error {
    // LogErrorf will output log and call fmt.Errorf to generate an error
    return log.LogErrorf(gxlog.LevelError, "%s", "another error")
}
```

