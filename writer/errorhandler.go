package writer

import (
	"fmt"
	"log"

	"github.com/gxlog/gxlog/iface"
)

// The ErrorHandler type is a function type used to handle errors.
// Do NOT call any method of the Writer or the Logger within the function,
// or it may deadlock.
type ErrorHandler func(bs []byte, record *iface.Record, err error)

// Report calls log.Output with the err.
func Report(_ []byte, _ *iface.Record, err error) {
	log.Output(1, fmt.Sprintln("log write error:", err))
}

// ReportDetails calls log.Output with the err and the bs.
func ReportDetails(bs []byte, _ *iface.Record, err error) {
	log.Output(1, fmt.Sprintf("log write error: %s, log: %s\n", err, bs))
}
