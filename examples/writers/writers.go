package main

import (
	"compress/flate"
	"fmt"
	"os"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/defaults"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/writer"
	"github.com/gxlog/gxlog/writer/file"
	"github.com/gxlog/gxlog/writer/socket/tcp"
	"github.com/gxlog/gxlog/writer/socket/unix"
	"github.com/gxlog/gxlog/writer/syslog"
)

var log = defaults.Logger()

func main() {
	// Only supported on systems that ANSI escape sequences are supported.
	defaults.Formatter().EnableColor()

	testWrappers()
	testSocketWriters()

	defaults.Formatter().DisableColor()

	defaults.Formatter().SetHeader(text.CompactHeader)
	testFileWriter()

	defaults.Formatter().SetHeader(text.SyslogHeader)
	testSyslogWriter()
}

func testWrappers() {
	// custom writer function
	fn := writer.Func(func(bs []byte, _ *gxlog.Record) {
		os.Stderr.Write(bs)
	})
	log.SetSlotWriter(gxlog.Slot0, fn)
	log.Info("a simple writer that just writes to os.Stderr")

	// use wrapper of io.Writer as another equivalent way
	wt := writer.Wrap(os.Stderr)
	log.SetSlotWriter(gxlog.Slot0, wt)
	log.Info("writer wrapper of os.Stderr")

	// asynchronous writer wrapper which uses a internal channel to buffer logs
	// when the channel is full, the Write method of the wrapper blocks
	// ATTENTION: some logs may NOT be output in asynchronous mode if os.Exit
	//   is called, panicking without recovery or ...
	async := writer.NewAsync(wt, 1024)
	// Close waits until all logs have been output.
	// It does NOT close the underlying writer.
	// To ignore all logs that have not been output, use Abort instead.
	defer async.Close()

	log.SetSlotWriter(gxlog.Slot0, async)
	log.Info("asynchronous writer wrapper")
}

func testSocketWriters() {
	// tcp socket writer
	// For performance and security, use a unix writer instead as long as the
	// system has support for unix domain socket.
	tcpWriter, err := tcp.Open(tcp.NewConfig(":9999"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tcpWriter.Close()

	log.SetSlotWriter(gxlog.Slot0, tcpWriter)

	// use "netcat localhost 9999" to watch logs
	// for i := 0; i < 1024; i++ {
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
	// }

	// shell expansion is NOT supported, so ~, $var and so on will not be expanded
	unixWriter, err := unix.Open(unix.NewConfig("/tmp/gxlog/unixdomain"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer unixWriter.Close()

	log.SetSlotWriter(gxlog.Slot0, unixWriter)

	// use "netcat -U /tmp/gxlog/unixdomain" to watch logs
	// for i := 0; i < 1024; i++ {
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
	// }
}

func testFileWriter() {
	// shell expansion is NOT supported, so ~, $var and so on will not be expanded
	wt, err := file.Open(file.NewConfig("/tmp/gxlog", "test").
		WithDateStyle(file.DateUnderscore).
		WithTimeStyle(file.TimeDot))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wt.Close()

	log.SetSlotWriter(gxlog.Slot0, wt)
	log.Info("this will be output to a file")

	wt.UpdateConfig(func(config file.Config) file.Config {
		// Do NOT call methods of the file writer, or it will deadlock.
		config.Ext = ".bin"
		// enable gzip compression
		config.GzipLevel = flate.DefaultCompression
		// enable AES encryption. the key must be hexadecimal encoded.
		config.AESKey = "70856575b161fbcca8fc12e1f70fc1c8"
		return config
	})
	log.Info("gzipped and encrypted")
}

func testSyslogWriter() {
	// connect to the local syslog server
	wt, err := syslog.Open(syslog.NewConfig("gxlog"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wt.Close()

	log.SetSlotWriter(gxlog.Slot0, wt)
	// NOTICE: the std syslog package will get the timestamp itself which is
	// a tiny bit later than Record.Time.
	log.Info("this will be output to syslog")

	// update level mapping
	// the severity of a level is left to be unchanged if it is not in the map
	wt.MapSeverity(map[gxlog.Level]syslog.Severity{
		gxlog.Info: syslog.SevErr,
	})
	log.Info("this will be severity err")
}
