package writer

import (
	"github.com/gxlog/gxlog"
)

type logData struct {
	Bytes  []byte
	Record *gxlog.Record
}

// An Async is a wrapper to the interface gxlog.Writer.
// All writers of gxlog.Writer it wraps switch into asynchronous mode.
//
// All methods of an Async are concurrency safe.
type Async struct {
	writer    gxlog.Writer
	chanData  chan logData
	chanClose chan struct{}
}

// NewAsync creates a new Async that wraps the writer. The writer must NOT be nil.
// The cap is the capacity of the internal channel of the Async and it must NOT
// be negative.
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

// Write implements the interface gxlog.Writer. It sends the bs and record to the
// internal channel. Another goroutine will receive them from the channel and
// call the underlying writer to output logs. If the channel is full, it blocks.
func (async *Async) Write(bs []byte, record *gxlog.Record) {
	async.chanData <- logData{Bytes: bs, Record: record}
}

// Close closes the internal channel and waits until all logs in the channel
// have been output. It does NOT close the underlying writer.
func (async *Async) Close() {
	close(async.chanClose)
	close(async.chanData)
	for data := range async.chanData {
		async.writer.Write(data.Bytes, data.Record)
	}
}

// Abort closes the internal channel and ignores all logs in the channel.
// It does NOT close the underlying writer.
func (async *Async) Abort() {
	close(async.chanClose)
	close(async.chanData)
}

// Len returns the len of the internal channel.
func (async *Async) Len() int {
	return len(async.chanData)
}

func (async *Async) serve() {
	for {
		select {
		case data := <-async.chanData:
			async.writer.Write(data.Bytes, data.Record)
		case <-async.chanClose:
			return
		}
	}
}
