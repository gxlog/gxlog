package json

import (
	"fmt"
)

const cCtrlCharN = 0x20

var gMapArray [cCtrlCharN]string

func init() {
	for i := 0; i < cCtrlCharN; i++ {
		gMapArray[i] = fmt.Sprintf(`\u00%02x`, i)
	}
}

func escape(buf []byte, str string) []byte {
	for i := 0; i < len(str); i++ {
		b := str[i]
		if b < cCtrlCharN {
			switch b {
			case '\n':
				buf = append(buf, `\n`...)
			case '\r':
				buf = append(buf, `\r`...)
			case '\t':
				buf = append(buf, `\t`...)
			default:
				buf = append(buf, gMapArray[b]...)
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
