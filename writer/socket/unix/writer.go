// Package unix implements a unix domain socket writer which implements the iface.Writer.
//
// The unix domain socket writer aims at log watching. For log transmission, use
// a syslog writer instead. With a unix domain socket writer, you can use netcat
// to receive and watch logs rather than the tail which is inconvenient because
// a new log file will be created when a log file reaches the max size.
package unix

import (
	"fmt"
	"os"

	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/writer/socket/internal/socket"
)

// A Writer implements the interface iface.Writer.
//
// All methods of a Writer are concurrency safe.
// A Writer MUST be created with Open.
type Writer struct {
	writer *socket.Writer
}

// Open creates a new Writer with the config.
func Open(config Config) (*Writer, error) {
	config.setDefaults()
	if !config.NoOverwrite {
		if err := checkAndRemove(config.Pathname); err != nil {
			return nil, fmt.Errorf("writer/socket/unix.Open: %v", err)
		}
	}
	writer, err := socket.Open("unix", config.Pathname)
	if err != nil {
		return nil, fmt.Errorf("writer/socket/unix.Open: %v", err)
	}
	if err := os.Chmod(config.Pathname, config.Perm); err != nil {
		writer.Close()
		return nil, fmt.Errorf("writer/socket/unix.Open: %v", err)
	}
	return &Writer{writer: writer}, nil
}

// Close closes the Writer.
func (writer *Writer) Close() error {
	if err := writer.writer.Close(); err != nil {
		return fmt.Errorf("writer/socket/unix.Close: %v", err)
	}
	return nil
}

// Write implements the interface Writer. It writes logs to unix domain sockets.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.writer.Write(bs, record)
}

func checkAndRemove(pathname string) error {
	if _, err := os.Stat(pathname); err != nil {
		return nil
	}
	return os.Remove(pathname)
}
