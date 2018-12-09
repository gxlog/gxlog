package writer

import (
	"github.com/gratonos/gxlog"
)

type logData struct {
	bs     []byte
	record *gxlog.Record
}

type Async struct {
	writer    gxlog.Writer
	chanData  chan logData
	chanClose chan struct{}
}

func NewAsync(wt gxlog.Writer, chanLen int) *Async {
	if wt == nil {
		panic("nil wt")
	}
	async := &Async{
		writer:    wt,
		chanData:  make(chan logData, chanLen),
		chanClose: make(chan struct{}),
	}
	go async.serve()
	return async
}

func (this *Async) Write(bs []byte, record *gxlog.Record) {
	this.chanData <- logData{bs: bs, record: record}
}

func (this *Async) Close() {
	close(this.chanClose)
	close(this.chanData)
	for data := range this.chanData {
		this.writer.Write(data.bs, data.record)
	}
}

func (this *Async) Abort() {
	close(this.chanClose)
	close(this.chanData)
}

func (this *Async) serve() {
	for {
		select {
		case data := <-this.chanData:
			this.writer.Write(data.bs, data.record)
		case <-this.chanClose:
			break
		}
	}
}
