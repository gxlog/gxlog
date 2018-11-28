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
	wt, err := socket.Open("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("tcp.Open: %v", err)
	}
	return &Writer{wt}, nil
}

func (this *Writer) Close() error {
	return this.writer.Close()
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.writer.Write(bs, record)
}
