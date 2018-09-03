package text

import "fmt"

type Formatter struct{}

func (*Formatter) Format(record *Record) []byte {
	return []byte(fmt.Sprint(record))
}
