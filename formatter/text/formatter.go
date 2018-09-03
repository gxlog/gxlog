package text

import "strconv"

type Formatter struct{}

func (*Formatter) Format(i int) []byte {
	return []byte(strconv.Itoa(i) + "\n")
}
