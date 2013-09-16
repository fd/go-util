package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

var StackLimit = 3
var StackContext = 3

type Error struct {
	err     error
	message string
	context []string
	stack   Stack
}

func Guard(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewFromPanic(r)
		}
	}()

	err = f()
	return
}

func NewFromPanic(r interface{}) *Error {
	if r == nil {
		return nil
	}

	if err, ok := r.(*Error); ok {
		return Annotate(err, "panic: %s", err.message)
	}

	if err, ok := r.(error); ok {
		return Annotate(err, "panic: %s", err)
	}

	if msg, ok := r.(string); ok {
		return New("panic: %s", msg)
	}

	return New("panic: %+v", r)
}

func New(message string, args ...interface{}) *Error {
	return &Error{
		message: fmt.Sprintf(message, args...),
		stack:   CaptureStack().Skip(2),
	}
}

func Annotate(err error, message string, args ...interface{}) *Error {
	if l, ok := err.(List); ok {
		err = l.Normalize()
	}

	if err == nil {
		return nil
	}

	return &Error{
		err:     err,
		message: fmt.Sprintf(message, args...),
		stack:   CaptureStack().Skip(2),
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
		stack := e.stack.Limit(StackLimit).String()
		stack = strings.Replace(stack, "\n", "\n    ", -1)
		fmt.Fprintf(&buf, "  location:\n    %s\n", stack)

		// stack := strings.TrimSpace(e.stack)
		// parts := strings.Split(stack, "\n")
		// if len(parts) > StackSize*2 {
		//   parts = parts[:StackSize*2]
		// }
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
		Stack:   e.stack.Limit(StackLimit).String(),
		Context: ctx,
	}

	return json.Marshal(json_err)
}
