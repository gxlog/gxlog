package tcp

import (
	"fmt"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/writer/socket/internal/socket"
)

type Writer struct {
	writer *socket.Writer
}

func Open(config *Config) (*Writer, error) {
	if config == nil {
		panic("nil config")
	}
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("tcp.Open: %v", err)
	}
	wt, err := socket.Open("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("tcp.Open: %v", err)
	}
	return &Writer{writer: wt}, nil
}

func (this *Writer) Close() error {
	if err := this.writer.Close(); err != nil {
		return fmt.Errorf("tcp.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.writer.Write(bs, record)
}
