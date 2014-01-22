// +build !heroku

package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05.000"

func (l *logger) format(level Level, s string) string {
	var (
		stdout = false
	)

	if f, ok := l.w.(*os.File); ok && f.Fd() == os.Stdout.Fd() {
		stdout = true
	}

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
	if !stdout {
		if l.namespace == "" {
			s = fmt.Sprintf("%s [%c] %s\n", now, code, s)
		} else {
			s = fmt.Sprintf("%s [%c] %s: %s\n", now, code, l.namespace, s)
		}
	} else {
		if l.namespace == "" {
			s = fmt.Sprintf("\x1B[33m%s [%c]\x1B[0m %s\n", now, code, s)
		} else {
			s = fmt.Sprintf("\x1B[33m%s [%c]\x1B[0m \x1B[34m%s:\x1B[0m %s\n", now, code, l.namespace, s)
		}
	}

	return s
}
