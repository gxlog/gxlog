package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type Formatter struct{}

func (*Formatter) Format(record *gxlog.Record) []byte {
	return []byte(fmt.Sprintln(record))
}
