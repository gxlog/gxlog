// Package tcp implements a tcp socket writer which implements the iface.Writer.
//
// The tcp socket writer aims at log watching. For log transmission, use a syslog
// writer instead. With a tcp socket writer, one can use netcat to receive logs
// rather than to tail a log file which is inconvenient because a new log file
// will be created when a log file reaches the max size.
//
// For performance and security, use a unix writer instead as long as the system
// has support for unix domain socket.
package tcp

import (
	"fmt"

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
	writer, err := socket.Open("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("writer/socket/tcp.Open: %v", err)
	}
	return &Writer{writer: writer}, nil
}

// Close closes the Writer.
func (writer *Writer) Close() error {
	if err := writer.writer.Close(); err != nil {
		return fmt.Errorf("writer/socket/tcp.Close: %v", err)
	}
	return nil
}

// Write implements the interface iface.Writer. It writes logs to tcp sockets.
func (writer *Writer) Write(bs []byte, record *iface.Record) {
	writer.writer.Write(bs, record)
}
