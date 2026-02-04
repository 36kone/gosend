package utils

import "fmt"

func ColorizeStatusCode(code int) string {
	var color string

	switch {
	case code >= 100 && code < 200:
		color = Gray
	case code >= 200 && code < 300:
		color = Green
	case code >= 300 && code < 400:
		color = Yellow
	case code >= 400 && code < 500:
		color = RedLight
	case code >= 500:
		color = RedStrong
	default:
		color = Reset
	}

	return fmt.Sprintf("%s%d%s", color, code, Reset)
}
