package socket

import (
	"fmt"
	"net"
	"sync"

	"github.com/gxlog/gxlog"
)

type Writer struct {
	listener net.Listener
	conns    map[int64]net.Conn
	id       int64
	wg       sync.WaitGroup

	lock sync.Mutex
}

func Open(network, addr string) (*Writer, error) {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return nil, fmt.Errorf("socket.Open: %v", err)
	}
	wt := &Writer{
		listener: listener,
		conns:    make(map[int64]net.Conn),
	}
	wt.wg.Add(1)
	go wt.serve()
	return wt, nil
}

func (writer *Writer) Close() error {
	if err := writer.listener.Close(); err != nil {
		return fmt.Errorf("socket.Close: %v", err)
	}

	writer.wg.Wait()

	writer.lock.Lock()
	defer writer.lock.Unlock()

	for id, conn := range writer.conns {
		conn.Close()
		delete(writer.conns, id)
	}

	return nil
}

func (writer *Writer) Write(bs []byte, record *gxlog.Record) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	for id, conn := range writer.conns {
		if _, err := conn.Write(bs); err != nil {
			conn.Close()
			delete(writer.conns, id)
		}
	}
}

func (writer *Writer) serve() {
	for {
		conn, err := writer.listener.Accept()
		if err != nil {
			break
		}

		writer.lock.Lock()

		id := writer.id
		writer.id++
		writer.conns[id] = conn

		writer.lock.Unlock()
	}
	writer.wg.Done()
}
