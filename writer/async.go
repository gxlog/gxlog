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

func (this *Async) Write(bs []byte, record *gxlog.Record) {
	this.chanData <- logData{Bytes: bs, Record: record}
}

func (this *Async) Close() {
	close(this.chanClose)
	close(this.chanData)
	for data := range this.chanData {
		this.writer.Write(data.Bytes, data.Record)
	}
}

func (this *Async) Abort() {
	close(this.chanClose)
	close(this.chanData)
}

func (this *Async) Len() int {
	return len(this.chanData)
}

func (this *Async) serve() {
	for {
		select {
		case data := <-this.chanData:
			this.writer.Write(data.Bytes, data.Record)
		case <-this.chanClose:
			break
		}
	}
}
