package gxlog

// And returns a function that is the logic AND of all the filters.
// It has the short circuit feature.
// There must be at least one filter and any filter must NOT be nil.
func And(filters ...Filter) Filter {
	if len(filters) == 0 {
		panic("gxlog.And: no filters")
	}
	return func(record *Record) bool {
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
// There must be at least one filter and any filter must NOT be nil.
func Or(filters ...Filter) Filter {
	if len(filters) == 0 {
		panic("gxlog.Or: no filters")
	}
	return func(record *Record) bool {
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
	return func(record *Record) bool {
		return !filter(record)
	}
}
