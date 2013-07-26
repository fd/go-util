package errors

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"text/tabwriter"
)

type Error struct {
	err     error
	message string
	context []string
	stack   string
}

func New(message string, args ...interface{}) *Error {
	stack := string(debug.Stack())
	parts := strings.SplitN(string(debug.Stack()), "\n", 3)
	stack = parts[2]

	return &Error{
		message: fmt.Sprintf(message, args...),
		stack:   stack,
	}
}

func Annotate(err error, message string, args ...interface{}) *Error {
	stack := string(debug.Stack())
	parts := strings.SplitN(string(debug.Stack()), "\n", 3)
	stack = parts[2]

	return &Error{
		err:     err,
		message: fmt.Sprintf(message, args...),
		stack:   stack,
	}
}

func (e *Error) AddContext(format string, args ...interface{}) {
	e.context = append(
		e.context,
		fmt.Sprintf(format, args...),
	)
}

func (e *Error) Error() string {
	var (
		buf bytes.Buffer
	)

	fmt.Fprintf(&buf, "error: %s\n", e.message)

	if len(e.context) > 0 {
		fmt.Fprintln(&buf, "  context:")

		sort.Strings(e.context)

		w := tabwriter.NewWriter(&buf, 8, 8, 1, ' ', 0)

		for _, pair := range e.context {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 1 {
				fmt.Fprintf(w, "    %s\t\n", parts[0])
			} else {
				fmt.Fprintf(w, "    %s\t= %s\n", parts[0], parts[1])
			}
			w.Flush()
		}
	}

	if len(e.stack) > 0 {
		fmt.Fprintln(&buf, "  location:")

		parts := strings.Split(strings.TrimSpace(e.stack), "\n")
		if len(parts) > 6 {
			parts = parts[:6]
		}

		for _, line := range parts {
			fmt.Fprintf(&buf, "    %s\n", line)
		}
	}

	if e.err != nil {
		err := e.err.Error()
		err = strings.TrimSpace(err)
		err = strings.Replace(err, "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "  %s\n", err)
	}

	return buf.String()
}

/*

error: Error message
  context:
    foo:   bar
    hello: world
  location:
    example/path.go:34
  error: Error message
    ...


*/
