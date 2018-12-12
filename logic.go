package gxlog

func And(filters ...Filter) Filter {
	return func(record *Record) bool {
		for _, filter := range filters {
			if !filter(record) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...Filter) Filter {
	return func(record *Record) bool {
		for _, filter := range filters {
			if filter(record) {
				return true
			}
		}
		return false
	}
}

func Not(filter Filter) Filter {
	return func(record *Record) bool {
		return !filter(record)
	}
}
