package socket

import (
	"fmt"
	"net"
	"sync"

	"github.com/gratonos/gxlog"
)

type Writer struct {
	config   Config
	listener net.Listener
	conns    map[int64]net.Conn
	id       int64
	wg       sync.WaitGroup
	lock     sync.Mutex
}

func Open(config *Config) (*Writer, error) {
	listener, err := net.Listen(config.Network, config.Bind)
	if err != nil {
		return nil, fmt.Errorf("socket.Open: %v", err)
	}
	wt := &Writer{
		config:   *config,
		listener: listener,
		conns:    make(map[int64]net.Conn),
	}
	wt.wg.Add(1)
	go wt.serve()
	return wt, nil
}

func (this *Writer) Close() error {
	if err := this.listener.Close(); err != nil {
		return fmt.Errorf("socket.Close: %v", err)
	}
	this.wg.Wait()
	return nil
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {
	this.lock.Lock()
	for id, conn := range this.conns {
		if _, err := conn.Write(bs); err != nil {
			conn.Close()
			delete(this.conns, id)
		}
	}
	this.lock.Unlock()
}

func (this *Writer) serve() {
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			break
		}
		this.lock.Lock()
		id := this.id
		this.id++
		this.conns[id] = conn
		this.lock.Unlock()
	}
	this.wg.Done()
}