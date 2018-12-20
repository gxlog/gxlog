package unix

import (
	"fmt"
	"os"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/writer/socket/internal/socket"
)

type Writer struct {
	writer *socket.Writer
}

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

func (this *Writer) Close() error {
	if err := this.writer.Close(); err != nil {
		return fmt.Errorf("writer/socket/unix.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.writer.Write(bs, record)
}

func checkAndRemove(pathname string) error {
	if _, err := os.Stat(pathname); err != nil {
		return nil
	}
	return os.Remove(pathname)
}
