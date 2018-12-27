package file

import (
	"compress/gzip"
	"io"
)

type gzipWriter struct {
	underlying io.WriteCloser
	writer     *gzip.Writer
}

func newGzipWriter(wt io.WriteCloser, level int) (io.WriteCloser, error) {
	writer, err := gzip.NewWriterLevel(wt, level)
	if err != nil {
		return wt, err
	}
	return &gzipWriter{
		underlying: wt,
		writer:     writer,
	}, nil
}

func (gz *gzipWriter) Close() error {
	if err := gz.writer.Close(); err != nil {
		return err
	}
	return gz.underlying.Close()

}

func (gz *gzipWriter) Write(bs []byte) (n int, err error) {
	n, err = gz.writer.Write(bs)
	if err == nil {
		err = gz.writer.Flush()
	}
	return n, err
}
