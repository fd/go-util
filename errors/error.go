package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"text/tabwriter"
)

var StackSize = 3

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
	stack = strings.Replace(stack, "\t", "  ", -1)

	return &Error{
		message: fmt.Sprintf(message, args...),
		stack:   stack,
	}
}

func Annotate(err error, message string, args ...interface{}) *Error {
	if l, ok := err.(List); ok {
		err = l.Normalize()
	}

	if err == nil {
		return nil
	}

	stack := string(debug.Stack())
	parts := strings.SplitN(string(debug.Stack()), "\n", 3)
	stack = parts[2]
	stack = strings.Replace(stack, "\t", "  ", -1)

	return &Error{
		err:     err,
		message: fmt.Sprintf(message, args...),
		stack:   stack,
	}
}

func (e *Error) AddContext(format string, args ...interface{}) {
	if e == nil {
		return
	}

	e.context = append(
		e.context,
		fmt.Sprintf(format, args...),
	)
}

func (e *Error) Error() string {
	if e == nil {
		return "(no error)"
	}

	var (
		buf            bytes.Buffer
		basic_child    bool
		extended_child bool
	)

	switch e.err.(type) {
	case *Error, List:
		extended_child = true
	default:
		if e.err != nil {
			basic_child = true
		}
	}

	fmt.Fprintf(&buf, "error: %s\n", e.message)

	if basic_child {
		fmt.Fprintf(&buf, "  message: %s\n", e.err)
	}

	if len(e.context) > 0 {
		fmt.Fprintln(&buf, "  context:")

		sort.Strings(e.context)

		w := tabwriter.NewWriter(&buf, 0, 8, 1, ' ', 0)

		for _, pair := range e.context {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 1 {
				fmt.Fprintf(w, "    %s\t\n", parts[0])
			} else {
				fmt.Fprintf(w, "    %s\t= %s\n", parts[0], parts[1])
			}
		}

		w.Flush()
	}

	if len(e.stack) > 0 {
		fmt.Fprintln(&buf, "  location:")

		stack := strings.TrimSpace(e.stack)
		parts := strings.Split(stack, "\n")
		if len(parts) > StackSize*2 {
			parts = parts[:StackSize*2]
		}

		for _, line := range parts {
			fmt.Fprintf(&buf, "    %s\n", line)
		}
	}

	if extended_child {
		err := e.err.Error()
		err = strings.TrimSpace(err)
		err = strings.Replace(err, "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "  %s\n", err)
	}

	return buf.String()
}

type json_Error struct {
	Message string                 `json:"message,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	Stack   string                 `json:"stack,omitempty"`
	Error   interface{}            `json:"error,omitempty"`
}

func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	var (
		ctx map[string]interface{}
	)

	if len(e.context) > 0 {
		ctx = make(map[string]interface{}, len(e.context))
		for _, pair := range e.context {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 1 {
				ctx[parts[0]] = true
			} else {
				ctx[parts[0]] = parts[1]
			}
		}
	}

	json_err := json_Error{
		Error:   e.err,
		Message: e.message,
		Stack:   e.stack,
		Context: ctx,
	}

	return json.Marshal(json_err)
}
