package gxlog

import "io"

type Formatter interface {
	Format(int) []byte
}

type Writer interface {
	io.Writer
}
