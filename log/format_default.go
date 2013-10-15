// +build !heroku

package log

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05.000"

func (l *logger) format(level Level, s string) string {
	// format the time of the message
	now := time.Now().UTC().Format(timeFormat)

	// clean up whitespace, indent subsequent lines
	s = strings.TrimSpace(strings.Replace(s, "\n", "\n  ", -1))

	// get the level code
	code := level_codes[level]
	if code == 0 {
		code = level_codes[UNKNOWN]
	}

	// format the message
	s = fmt.Sprintf("%s [%c] %s\n", now, code, s)

	return s
}
