package writer

import (
	"fmt"
	"log"

	"github.com/gxlog/gxlog/iface"
)

const callDepthOffset = 6

// The ErrorHandler type is a function type used to handle errors.
// Do NOT call any method of the Writer or the Logger within the function,
// or it may deadlock.
type ErrorHandler func(bs []byte, record *iface.Record, err error)

// Report calls log.Output with the err.
func Report(_ []byte, _ *iface.Record, err error) {
	log.Output(callDepthOffset, fmt.Sprintln("log error:", err))
}

// ReportDetails calls log.Output with the err and the bs.
func ReportDetails(bs []byte, _ *iface.Record, err error) {
	log.Output(callDepthOffset, fmt.Sprintf("log error: %s, log: %s", err, bs))
}
