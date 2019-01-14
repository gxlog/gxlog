package logger

import (
	"github.com/gxlog/gxlog/iface"
)

// And returns a function that is the logic AND of all the filters.
// It has the short circuit feature.
// Any filter must NOT be nil.
func And(filter Filter, filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		if !filter(record) {
			return false
		}
		for _, filter := range filters {
			if !filter(record) {
				return false
			}
		}
		return true
	}
}

// Or returns a function that is the logic OR of all the filters.
// It has the short circuit feature.
// Any filter must NOT be nil.
func Or(filter Filter, filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		if filter(record) {
			return true
		}
		for _, filter := range filters {
			if filter(record) {
				return true
			}
		}
		return false
	}
}

// Not returns a function that is the logic NOT of the filter.
// The filter must NOT be nil.
func Not(filter Filter) Filter {
	return func(record *iface.Record) bool {
		return !filter(record)
	}
}
