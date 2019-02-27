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
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
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
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
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
	log.Info("this will be output to syslog")

	// update level mapping
	// The severity of a level is left to be unchanged if it is not in the map.
	wt.MapSeverities(map[iface.Level]syslog.Severity{
		iface.Info: syslog.SevErr,
	})
	log.Info("this will be severity err")
}
