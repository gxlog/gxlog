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

## Architecture ##

## Features Preview ##

- logger
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
  - multi-slot
    - slots management
    - slots level
    - slots filter
  - formatters
    - formatter function wrapper
    - null formatter
    - text formatter
      - custom header
      - custom property and format of fields
      - colorization
      - custom color mapping
    - json formatter
      - custom property of fields
      - custom omission of fields
      - custom omission of empty fields
  - writer
    - writer function wrapper
    - io.Writer wrapper
    - asynchronous
    - null writer
    - file writer
      - custom file max size
      - custom file naming style
      - file deletion check
      - new directory each day
      - gzip compression
      - AES encryption
      - error reporting
    - syslog writer
      - custom mapping from level to severity
      - error reporting
    - tcp socket writer
    - unix domain socket writer

## Getting Started ##

