// Package unix implements a unix domain socket writer which implements the iface.Writer.
//
// The unix domain socket writer aims at log watching. For log transmission, use
// a syslog writer instead. With a unix domain socket writer, one can use netcat
// to receive logs rather than to tail a log file which is inconvenient because
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
//
// A Writer must be created with Open.
type Writer struct {
	writer *socket.Writer
}

// Open creates a new Writer with the config. The config must NOT be nil.
func Open(config *Config) (*Writer, error) {
	if config.Overwrite {
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

// Write implements the interface iface.Writer. It writes logs to unix domain sockets.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.writer.Write(bs, record)
}

func checkAndRemove(pathname string) error {
	if _, err := os.Stat(pathname); err != nil {
		return nil
	}
	return os.Remove(pathname)
}
