package gxlog

// And returns a function that is the logic AND of all the filters.
// Nil filters will be ignored. It has the short circuit feature.
// If the count of non-nil filters is zero, the returned function
// will always return true.
func And(filters ...Filter) Filter {
	return func(record *Record) bool {
		for _, filter := range filters {
			if filter != nil && !filter(record) {
				return false
			}
		}
		return true
	}
}

// Or returns a function that is the logic OR of all the filters.
// Nil filters will be ignored. It has the short circuit feature.
// If the count of non-nil filters is zero, the returned function
// will always return false.
func Or(filters ...Filter) Filter {
	return func(record *Record) bool {
		for _, filter := range filters {
			if filter != nil && filter(record) {
				return true
			}
		}
		return false
	}
}

// Not returns a function that is the logic NOT of the filter.
// If the filter is nil, the returned function will always return false.
func Not(filter Filter) Filter {
	return func(record *Record) bool {
		if filter != nil {
			return !filter(record)
		}
		return false
	}
}
