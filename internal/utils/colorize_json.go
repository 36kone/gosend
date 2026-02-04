package utils

import (
	"bytes"
)

func ColorizeJSON(jsonStr string) string {
	var buf bytes.Buffer
	inQuotes := false

	for i := 0; i < len(jsonStr); i++ {
		c := jsonStr[i]

		switch c {
		case '"':
			inQuotes = !inQuotes
			buf.WriteString(Green)
			buf.WriteByte(c)
			continue
		case '{', '}', '[', ']':
			buf.WriteString(Blue)
		case ':':
			buf.WriteString(Yellow)
			buf.WriteByte(c)
			continue
		case ',':
			buf.WriteString(Blue)
			buf.WriteByte(c)
			continue
		}

		if inQuotes {
			buf.WriteByte(c)
			if i+1 < len(jsonStr) && jsonStr[i+1] == '"' {
				buf.WriteString(Reset)
			}
		} else {
			buf.WriteString(Reset)
			buf.WriteByte(c)
		}
	}

	buf.WriteString(Reset)
	return buf.String()
}
