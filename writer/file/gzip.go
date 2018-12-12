package file

import (
	"compress/gzip"
	"io"
)

type gzipWriter struct {
	underlying io.WriteCloser
	writer     *gzip.Writer
}

func newGzipWriter(wt io.WriteCloser, level int) *gzipWriter {
	// the level is ensured to be valid, then NewWriterLevel will not fail
	writer, _ := gzip.NewWriterLevel(wt, level)
	return &gzipWriter{
		underlying: wt,
		writer:     writer,
	}
}

func (this *gzipWriter) Close() error {
	if err := this.writer.Close(); err != nil {
		return err
	}
	return this.underlying.Close()

}

func (this *gzipWriter) Write(bs []byte) (n int, err error) {
	n, err = this.writer.Write(bs)
	if err == nil {
		err = this.writer.Flush()
	}
	return n, err
}
