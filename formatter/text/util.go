package text

func lastSegments(str string, n int, sep byte) string {
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == sep {
			n--
			if n == 0 {
				return str[i+1:]
			}
		}
	}
	return str
}
