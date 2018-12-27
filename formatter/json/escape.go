package json

import (
	"fmt"
)

const ctrlCharCount = 0x20

var escapeMapArray [ctrlCharCount]string

func init() {
	for i := 0; i < ctrlCharCount; i++ {
		escapeMapArray[i] = fmt.Sprintf(`\u00%02x`, i)
	}
}

func escape(buf []byte, str string) []byte {
	for i := 0; i < len(str); i++ {
		b := str[i]
		if b < ctrlCharCount {
			switch b {
			case '\n':
				buf = append(buf, `\n`...)
			case '\r':
				buf = append(buf, `\r`...)
			case '\t':
				buf = append(buf, `\t`...)
			default:
				buf = append(buf, escapeMapArray[b]...)
			}
		} else {
			switch b {
			case '"':
				buf = append(buf, `\"`...)
			case '\\':
				buf = append(buf, `\\`...)
			case '\u007f': // DEL
				// noop
			default:
				buf = append(buf, b)
			}
		}
	}
	return buf
}
