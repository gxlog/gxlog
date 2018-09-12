package null

import "github.com/gratonos/gxlog"

type Formatter struct{}

func (Formatter) Format(*gxlog.Record) []byte {
	return nil
}
