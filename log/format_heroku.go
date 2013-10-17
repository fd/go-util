// +build heroku

package log

import (
	"fmt"
	"strings"
)

func (l *logger) format(level Level, s string) string {
	// clean up whitespace, indent subsequent lines
	s = strings.TrimSpace(strings.Replace(s, "\n", "\n  ", -1))

	// get the level code
	code := level_codes[level]
	if code == 0 {
		code = level_codes[UNKNOWN]
	}

	// format the message
	if l.namespace == "" {
		s = fmt.Sprintf("[%c] %s\n", code, s)
	} else {
		s = fmt.Sprintf("[%c] %s: %s\n", code, l.namespace, s)
	}

	return s
}
