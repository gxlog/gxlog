package tcp

import (
	"fmt"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/writer/socket/internal/socket"
)

type Writer struct {
	writer *socket.Writer
}

func Open(config *Config) (*Writer, error) {
	writer, err := socket.Open("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("writer/socket/tcp.Open: %v", err)
	}
	return &Writer{writer: writer}, nil
}

func (this *Writer) Close() error {
	if err := this.writer.Close(); err != nil {
		return fmt.Errorf("writer/socket/tcp.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.writer.Write(bs, record)
}
