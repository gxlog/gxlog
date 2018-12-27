package writer

import (
	"github.com/gxlog/gxlog"
)

type logData struct {
	Bytes  []byte
	Record *gxlog.Record
}

type Async struct {
	writer    gxlog.Writer
	chanData  chan logData
	chanClose chan struct{}
}

func NewAsync(writer gxlog.Writer, cap int) *Async {
	if writer == nil {
		panic("writer.NewAsync: nil writer")
	}
	async := &Async{
		writer:    writer,
		chanData:  make(chan logData, cap),
		chanClose: make(chan struct{}),
	}
	go async.serve()
	return async
}

func (async *Async) Write(bs []byte, record *gxlog.Record) {
	async.chanData <- logData{Bytes: bs, Record: record}
}

func (async *Async) Close() {
	close(async.chanClose)
	close(async.chanData)
	for data := range async.chanData {
		async.writer.Write(data.Bytes, data.Record)
	}
}

func (async *Async) Abort() {
	close(async.chanClose)
	close(async.chanData)
}

func (async *Async) Len() int {
	return len(async.chanData)
}

func (async *Async) serve() {
	for {
		select {
		case data := <-async.chanData:
			async.writer.Write(data.Bytes, data.Record)
		case <-async.chanClose:
			break
		}
	}
}
